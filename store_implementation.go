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

var _ StoreInterface = (*storeImplementation)(nil) // verify it extends the interface

// storeImplementation is the concrete implementation of the StoreInterface.
// It provides database operations for posts, taxonomies, terms, and term relations.
type storeImplementation struct {
	postTableName         string
	taxonomyTableName     string
	termTableName         string
	termRelationTableName string
	db                    *sql.DB
	dbDriverName          string
	timeoutSeconds        int64
	automigrateEnabled    bool
	debugEnabled          bool

	versioningEnabled bool
	versioningStore   versionstore.StoreInterface

	taxonomyEnabled bool
}

// migrateSlugColumn adds the slug column if it doesn't exist (for existing installations)
// TODO: Remove this function after May 2027 (1 year from implementation)
func (store *storeImplementation) migrateSlugColumn() error {
	checkSQL, checkParams, err := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.postTableName).
		TableColumnExists(store.postTableName, COLUMN_SLUG)
	if err != nil {
		return err
	}

	var exists bool
	err = store.db.QueryRow(checkSQL, checkParams...).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		alterSQL, err := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
			Table(store.postTableName).
			TableColumnAdd(COLUMN_SLUG, sb.Column{
				Name:   COLUMN_SLUG,
				Type:   sb.COLUMN_TYPE_STRING,
				Length: 255,
			})
		if err != nil {
			return err
		}

		_, err = store.db.Exec(alterSQL)
		if err != nil {
			return err
		}
	}

	return nil
}

// AutoMigrate creates the necessary database tables for the blog store.
// It creates the post table, and optionally taxonomy tables if enabled.
// Also initializes the versioning store if versioning is enabled.
func (store *storeImplementation) AutoMigrate() error {
	// Create main post table
	sql, err := store.sqlCreatePostTable()
	if err != nil {
		return err
	}

	_, err = store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	// TODO: Remove this migration logic after May 2027 (1 year from implementation)
	// This allows existing installations to auto-migrate the slug column
	err = store.migrateSlugColumn()
	if err != nil {
		log.Println(err)
		return err
	}

	// Create taxonomy tables only if enabled
	if store.taxonomyEnabled {
		// Create taxonomy table
		sql, err = store.sqlCreateTaxonomyTable()
		if err != nil {
			return err
		}

		_, err = store.db.Exec(sql)
		if err != nil {
			log.Println(err)
			return err
		}

		// Create term table
		sql, err = store.sqlCreateTermTable()
		if err != nil {
			return err
		}

		_, err = store.db.Exec(sql)
		if err != nil {
			log.Println(err)
			return err
		}

		// Create term relation table
		sql, err = store.sqlCreateTermRelationTable()
		if err != nil {
			return err
		}

		_, err = store.db.Exec(sql)
		if err != nil {
			log.Println(err)
			return err
		}
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

// VersioningEnabled returns true if versioning is enabled for this store.
func (st *storeImplementation) VersioningEnabled() bool {
	return st.versioningEnabled
}

// TaxonomyEnabled returns true if taxonomy features are enabled for this store.
func (st *storeImplementation) TaxonomyEnabled() bool {
	return st.taxonomyEnabled
}

// EnableDebug enables or disables debug logging for SQL queries.
func (st *storeImplementation) EnableDebug(debug bool) StoreInterface {
	st.debugEnabled = debug
	return st
}

// PostCreate inserts a new post into the database.
// It sets the created_at and updated_at timestamps automatically.
// Also tracks the creation in the versioning store if versioning is enabled.
func (store *storeImplementation) PostCreate(ctx context.Context, post PostInterface) error {
	post.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	post.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	data := post.GetData()

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
	if err := store.versioningTrackEntity(ctx, VERSIONING_TYPE_POST, post.GetID(), post); err != nil {
		return err
	}

	return nil
}

// PostCount returns the total number of posts matching the given query options.
func (store *storeImplementation) PostCount(ctx context.Context, options PostQueryOptions) (int64, error) {
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

// PostTrash moves a post to trash by setting its status to POST_STATUS_TRASH.
func (store *storeImplementation) PostTrash(ctx context.Context, post PostInterface) error {
	post.SetStatus(POST_STATUS_TRASH)

	return store.PostUpdate(ctx, post)
}

// PostDelete permanently removes a post from the database.
func (store *storeImplementation) PostDelete(ctx context.Context, post PostInterface) error {
	if post == nil {
		return errors.New("post is nil")
	}

	return store.PostDeleteByID(ctx, post.GetID())
}

// PostDeleteByID permanently removes a post by its ID.
func (store *storeImplementation) PostDeleteByID(ctx context.Context, id string) error {
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

// PostFindByID retrieves a post by its ID.
// Supports both full IDs and shortened IDs with automatic unshortening.
func (store *storeImplementation) PostFindByID(ctx context.Context, id string) (PostInterface, error) {
	if id == "" {
		return nil, errors.New("post id is empty")
	}

	// Normalize ID
	normalizedID := NormalizeID(id)

	// Try direct lookup first
	list, err := store.PostList(ctx, PostQueryOptions{
		ID:    normalizedID,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	// If not found and ID looks shortened, try unshortening
	if IsShortID(normalizedID) {
		unshortened, err := UnshortenID(normalizedID)
		if err == nil && unshortened != normalizedID {
			list, err = store.PostList(ctx, PostQueryOptions{
				ID:    unshortened,
				Limit: 1,
			})

			if err != nil {
				return nil, err
			}

			if len(list) > 0 {
				return list[0], nil
			}
		}
	}

	return nil, nil
}

// PostFindPrevious finds the post created immediately before the given post.
func (st *storeImplementation) PostFindPrevious(post PostInterface) (PostInterface, error) {
	list, err := st.PostList(context.Background(), PostQueryOptions{
		CreatedAtLessThan: post.GetCreatedAtCarbon().ToDateTimeString(),
		Limit:             1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// PostFindNext finds the post created immediately after the given post.
func (st *storeImplementation) PostFindNext(post PostInterface) (PostInterface, error) {
	list, err := st.PostList(context.Background(), PostQueryOptions{
		CreatedAtGreaterThan: post.GetCreatedAtCarbon().ToDateTimeString(),
		Limit:                1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// PostList retrieves a list of posts matching the given query options.
func (st *storeImplementation) PostList(ctx context.Context, options PostQueryOptions) ([]PostInterface, error) {
	q := st.postQuery(options)

	sqlStr, sqlParams, errSql := q.Select().
		Prepared(true).
		ToSQL()

	if errSql != nil {
		log.Println(errSql)
		return []PostInterface{}, errSql
	}

	if st.debugEnabled {
		log.Println(sqlStr)
		log.Println(sqlParams)
	}

	modelMaps, err := database.SelectToMapString(
		database.NewQueryableContext(ctx, st.db),
		sqlStr,
		sqlParams...,
	)
	if err != nil {
		return []PostInterface{}, err
	}

	list := []PostInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewPostFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

// PostSoftDelete marks a post as deleted by setting the soft_deleted_at timestamp.
func (st *storeImplementation) PostSoftDelete(ctx context.Context, post PostInterface) error {
	if post == nil {
		return errors.New("post is nil")
	}

	post.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return st.PostUpdate(ctx, post)
}

// PostSoftDeleteByID marks a post as deleted by its ID.
func (st *storeImplementation) PostSoftDeleteByID(ctx context.Context, id string) error {
	post, err := st.PostFindByID(ctx, id)

	if err != nil {
		return err
	}

	return st.PostSoftDelete(ctx, post)
}

// PostUpdate updates an existing post in the database.
// Only changed fields are updated. Also tracks the update in the versioning store if enabled.
func (st *storeImplementation) PostUpdate(ctx context.Context, post PostInterface) error {
	if post == nil {
		return errors.New("order is nil")
	}

	// post.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := post.GetDataChanged()

	delete(dataChanged, "id")   // ID is not updatable
	delete(dataChanged, "hash") // Hash is not updatable
	delete(dataChanged, "data") // Data is not updatable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(st.dbDriverName).
		Update(st.postTableName).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(post.GetID())).
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
	if err2 := st.versioningTrackEntity(ctx, VERSIONING_TYPE_POST, post.GetID(), post); err2 != nil {
		return err2
	}

	return err
}

// postQuery builds a goqu SelectDataset for querying posts based on options.
func (st *storeImplementation) postQuery(options PostQueryOptions) *goqu.SelectDataset {
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

	// Handle MetaContains filtering - checks if meta JSON contains key-value pair
	if len(options.MetaContains) > 0 {
		for key, value := range options.MetaContains {
			// Use JSON extraction for cross-database compatibility
			// SQLite: json_extract(metas, '$.key')
			// MySQL/PostgreSQL: metas->>'$.key' or json_extract(metas, '$.key')
			var jsonExpr goqu.Expression
			switch st.dbDriverName {
			case "sqlite3", "sqlite":
				// SQLite json_extract syntax
				jsonExpr = goqu.L("json_extract("+COLUMN_METAS+", ?)", "$."+key).Eq(value)
			case "mysql":
				// MySQL JSON_EXTRACT syntax
				jsonExpr = goqu.L("JSON_UNQUOTE(JSON_EXTRACT("+COLUMN_METAS+", ?))", "$."+key).Eq(value)
			default:
				// PostgreSQL and others - use generic JSON extraction
				jsonExpr = goqu.L("("+COLUMN_METAS+"->>?)::text", key).Eq(value)
			}
			q = q.Where(jsonExpr)
		}
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
