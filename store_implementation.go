package blogstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dracory/versionstore"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

var _ StoreInterface = (*store)(nil) // verify it extends the interface

type store struct {
	postTableName      string
	db                 *sql.DB
	dbDriverName       string
	timeoutSeconds     int64
	automigrateEnabled bool
	debugEnabled       bool

	versioningEnabled bool
	versioningStore   versionstore.StoreInterface
}

// AutoMigrate auto migrate
func (store *store) AutoMigrate() error {
	sql := store.sqlCreateTable()

	_, err := store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	if store.versioningEnabled {
		if store.versioningStore == nil {
			return errors.New("versioning store is nil")
		}
		if err := store.versioningStore.AutoMigrate(); err != nil {
			return err
		}
	}

	return nil
}

func (st *store) VersioningEnabled() bool {
	return st.versioningEnabled
}

// EnableDebug - enables the debug option
func (st *store) EnableDebug(debug bool) StoreInterface {
	st.debugEnabled = debug
	return st
}

func (store *store) PostCreate(ctx context.Context, post *Post) error {
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

	_, err := store.db.ExecContext(ctx, sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	post.MarkAsNotDirty()
	if err := store.versioningTrackEntity(ctx, VERSIONING_TYPE_POST, post.ID(), post); err != nil {
		return err
	}

	return nil
}

func (store *store) PostCount(ctx context.Context, options PostQueryOptions) (int64, error) {
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

	mapped, err := database.SelectToMapString(
		database.NewQueryableContext(ctx, store.db),
		sqlStr,
		params...,
	)

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

func (store *store) PostTrash(ctx context.Context, post *Post) error {
	post.SetStatus(POST_STATUS_TRASH)

	return store.PostUpdate(ctx, post)
}

func (store *store) PostDelete(ctx context.Context, post *Post) error {
	if post == nil {
		return errors.New("post is nil")
	}

	return store.PostDeleteByID(ctx, post.ID())
}

func (store *store) PostDeleteByID(ctx context.Context, id string) error {
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

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	return err
}

func (store *store) PostFindByID(ctx context.Context, id string) (*Post, error) {
	if id == "" {
		return nil, errors.New("post id is empty")
	}

	list, err := store.PostList(ctx, PostQueryOptions{
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

func (st *store) PostFindPrevious(post Post) (*Post, error) {
	list, err := st.PostList(context.Background(), PostQueryOptions{
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

func (st *store) PostFindNext(post Post) (*Post, error) {
	list, err := st.PostList(context.Background(), PostQueryOptions{
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

func (st *store) PostList(ctx context.Context, options PostQueryOptions) ([]Post, error) {
	q := st.postQuery(options)

	sqlStr, sqlParams, errSql := q.Select().
		Prepared(true).
		ToSQL()

	if errSql != nil {
		log.Println(errSql)
		return []Post{}, errSql
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(
		database.NewQueryableContext(ctx, st.db),
		sqlStr,
		sqlParams...,
	)
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

func (st *store) PostSoftDelete(ctx context.Context, post *Post) error {
	if post == nil {
		return errors.New("post is nil")
	}

	post.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return st.PostUpdate(ctx, post)
}

func (st *store) PostSoftDeleteByID(ctx context.Context, id string) error {
	post, err := st.PostFindByID(ctx, id)

	if err != nil {
		return err
	}

	return st.PostSoftDelete(ctx, post)
}

func (st *store) PostUpdate(ctx context.Context, post *Post) error {
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

	sqlStr, params, errSql := goqu.Dialect(st.dbDriverName).
		Update(st.postTableName).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(post.ID())).
		Prepared(true).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.ExecContext(ctx, sqlStr, params...)

	post.MarkAsNotDirty()
	if err2 := st.versioningTrackEntity(ctx, VERSIONING_TYPE_POST, post.ID(), post); err2 != nil {
		return err2
	}

	return err
}

func (st *store) postQuery(options PostQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(st.dbDriverName).
		From(st.postTableName)

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
				// Search Title
				goqu.C(COLUMN_TITLE).ILike("%"+options.Search+"%"),
				// Search Body Content
				goqu.C(COLUMN_CONTENT).ILike("%"+options.Search+"%"),
				// Search ID
				goqu.C(COLUMN_ID).Eq(options.Search),
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
	}

	if options.WithDeleted {
		return q
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted)
}
