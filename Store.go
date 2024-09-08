package blogstore

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

type Store struct {
	postTableName      string
	db                 *sql.DB
	dbDriverName       string
	timeoutSeconds     int64
	automigrateEnabled bool
	debugEnabled       bool
}

// AutoMigrate auto migrate
func (store *Store) AutoMigrate() error {
	sql := store.sqlCreateTable()

	_, err := store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) *Store {
	st.debugEnabled = debug
	return st
}

func (store *Store) PostCreate(post *Post) error {
	post.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	post.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	data := post.Data()

	sqlStr, sqlParams, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.postTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	post.MarkAsNotDirty()

	return nil
}

func (store *Store) PostCount(options PostQueryOptions) (int64, error) {
	options.CountOnly = true
	q := store.postQuery(options)

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	mapped, err := db.SelectToMapString(sqlStr, params...)
	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err

	}

	return i, nil
}

func (store *Store) PostTrash(post *Post) error {
	post.SetStatus(POST_STATUS_TRASH)

	return store.PostUpdate(post)
}

func (store *Store) PostDelete(post *Post) error {
	if post == nil {
		return errors.New("post is nil")
	}

	return store.PostDeleteByID(post.ID())
}

func (store *Store) PostDeleteByID(id string) error {
	if id == "" {
		return errors.New("post id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.postTableName).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		Prepared(true).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

func (store *Store) PostFindByID(id string) (*Post, error) {
	if id == "" {
		return nil, errors.New("post id is empty")
	}

	list, err := store.PostList(PostQueryOptions{
		ID:    id,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

func (store *Store) PostFindPrevious(post Post) (*Post, error) {
	list, err := store.PostList(PostQueryOptions{
		CreatedAtLessThan: post.CreatedAtCarbon().ToDateTimeString(),
		Limit:             1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

func (store *Store) PostFindNext(post Post) (*Post, error) {
	list, err := store.PostList(PostQueryOptions{
		CreatedAtGreaterThan: post.CreatedAtCarbon().ToDateTimeString(),
		Limit:                1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

func (store *Store) PostList(options PostQueryOptions) ([]Post, error) {
	q := store.postQuery(options)

	sqlStr, sqlParams, errSql := q.Select().
		Prepared(true).
		ToSQL()

	if errSql != nil {
		log.Println(errSql)
		return []Post{}, errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr, sqlParams...)
	if err != nil {
		return []Post{}, err
	}

	list := []Post{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewPostFromExistingData(modelMap)
		list = append(list, *model)
	})

	return list, nil
}

func (store *Store) PostSoftDelete(post *Post) error {
	if post == nil {
		return errors.New("post is nil")
	}

	post.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.PostUpdate(post)
}

func (store *Store) PostSoftDeleteByID(id string) error {
	post, err := store.PostFindByID(id)

	if err != nil {
		return err
	}

	return store.PostSoftDelete(post)
}

func (store *Store) PostUpdate(post *Post) error {
	if post == nil {
		return errors.New("order is nil")
	}

	// post.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := post.DataChanged()

	delete(dataChanged, "id")   // ID is not updatable
	delete(dataChanged, "hash") // Hash is not updatable
	delete(dataChanged, "data") // Data is not updatable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.postTableName).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(post.ID())).
		Prepared(true).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	post.MarkAsNotDirty()

	return err
}

func (store *Store) postQuery(options PostQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).
		From(store.postTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if len(options.IDIn) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
	}

	if options.Search != "" {
		q = q.Where(
			goqu.Or(
				goqu.C(COLUMN_TITLE).Like("%"+options.Search+"%"),
				goqu.C(COLUMN_CONTENT).Like("%"+options.Search+"%"),
				goqu.C(COLUMN_ID).Like(options.Search),
			),
		)
	}

	if options.CreatedAtGreaterThan != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Gt(options.CreatedAtGreaterThan))
	}

	if options.CreatedAtLessThan != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Lt(options.CreatedAtLessThan))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	sortOrder := "desc"
	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.OrderBy != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy).Desc())
		}
	}

	if !options.WithDeleted {
		//q = q.Where(goqu.C("status").Neq(POST_STATUS_DELETED))
		q = q.Where(goqu.C(COLUMN_DELETED_AT).Eq(sb.NULL_DATETIME))
	}

	return q
}
