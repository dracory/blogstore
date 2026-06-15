package blogstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/dracory/neat"
	"github.com/dracory/neat/database/schema/constants"
	"github.com/dracory/versionstore"
	"github.com/dromara/carbon/v2"
	"log"
	"time"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
)

// StoreInterface defines the complete interface for blog post storage operations,
// including post management, taxonomy/term handling, and optional versioning support.
type StoreInterface interface {
	// GetPostTableName returns the post table name
	GetPostTableName() string
	// SetPostTableName sets the post table name
	SetPostTableName(tableName string)

	// GetTaxonomyTableName returns the taxonomy table name
	GetTaxonomyTableName() string
	// SetTaxonomyTableName sets the taxonomy table name
	SetTaxonomyTableName(tableName string)

	// GetTermTableName returns the term table name
	GetTermTableName() string
	// SetTermTableName sets the term table name
	SetTermTableName(tableName string)

	// GetTermRelationTableName returns the term relation table name
	GetTermRelationTableName() string
	// SetTermRelationTableName sets the term relation table name
	SetTermRelationTableName(tableName string)

	// MigrateDown drops the blog store tables
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error
	// MigrateUp creates the blog store tables
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error

	// EnableDebug toggles debug mode logging for database operations.
	// Returns the StoreInterface to allow method chaining.
	EnableDebug(debug bool) StoreInterface

	// VersioningEnabled returns true if versioning support is enabled for this store.
	VersioningEnabled() bool

	// TaxonomyEnabled returns true if taxonomy support is enabled for this store.
	TaxonomyEnabled() bool

	// PostCount returns the total number of posts matching the provided query options.
	// Uses PostQueryOptions to filter by status, type, or other criteria.
	PostCount(ctx context.Context, options PostQueryOptions) (int64, error)

	// PostCreate inserts a new post into the store.
	// Returns an error if the post cannot be created (e.g., duplicate ID or validation failure).
	PostCreate(ctx context.Context, post PostInterface) error

	// PostDelete permanently removes a post from the store.
	// This is a hard delete that cannot be undone.
	PostDelete(ctx context.Context, post PostInterface) error

	// PostDeleteByID permanently removes a post by its ID.
	// Returns an error if the post does not exist.
	PostDeleteByID(ctx context.Context, postID string) error

	// PostFindByID retrieves a post by its unique identifier.
	// Returns the post and nil error on success, or nil and an error if not found.
	PostFindByID(ctx context.Context, id string) (PostInterface, error)

	// PostFindBySlug retrieves a post by its slug.
	// Returns the post and nil error on success, or nil and an error if not found.
	PostFindBySlug(ctx context.Context, slug string) (PostInterface, error)

	// PostFindByOldSlug retrieves a post by its old slug (for redirect handling).
	// Returns the post and nil error on success, or nil and an error if not found.
	PostFindByOldSlug(ctx context.Context, oldSlug string) (PostInterface, error)

	// PostList retrieves a list of posts matching the provided query options.
	// Supports pagination, sorting, and filtering through PostQueryOptions.
	PostList(ctx context.Context, options PostQueryOptions) ([]PostInterface, error)

	// PostSoftDelete marks a post as deleted without removing it from the database.
	// The post can be restored later. Requires versioning to be enabled.
	PostSoftDelete(ctx context.Context, post PostInterface) error

	// PostSoftDeleteByID marks a post as deleted by ID without removing it from the database.
	// Returns an error if the post does not exist.
	PostSoftDeleteByID(ctx context.Context, postID string) error

	// PostTrash moves a post to the trash status.
	// Trashed posts are not visible in normal queries but can be restored.
	PostTrash(ctx context.Context, post PostInterface) error

	// PostUpdate modifies an existing post in the store.
	// Returns an error if the post does not exist or validation fails.
	PostUpdate(ctx context.Context, post PostInterface) error

	// Versioning methods manage historical versions of posts.

	// VersioningCreate saves a new version record for a post.
	VersioningCreate(ctx context.Context, versioning VersioningInterface) error

	// VersioningDelete permanently removes a versioning record.
	VersioningDelete(ctx context.Context, versioning VersioningInterface) error

	// VersioningDeleteByID permanently removes a versioning record by ID.
	VersioningDeleteByID(ctx context.Context, id string) error

	// VersioningFindByID retrieves a specific version record by ID.
	VersioningFindByID(ctx context.Context, versioningID string) (VersioningInterface, error)

	// VersioningList retrieves version records matching the provided query.
	VersioningList(ctx context.Context, query VersioningQueryInterface) ([]VersioningInterface, error)

	// VersioningSoftDelete marks a version record as deleted without permanent removal.
	VersioningSoftDelete(ctx context.Context, versioning VersioningInterface) error

	// VersioningSoftDeleteByID marks a version record as deleted by ID.
	VersioningSoftDeleteByID(ctx context.Context, id string) error

	// VersioningUpdate modifies an existing version record.
	VersioningUpdate(ctx context.Context, versioning VersioningInterface) error

	// Taxonomy methods manage classification systems (e.g., categories, tags).

	// TaxonomyCount returns the number of taxonomies matching the query options.
	TaxonomyCount(ctx context.Context, options TaxonomyQueryOptions) (int64, error)

	// TaxonomyCreate inserts a new taxonomy into the store.
	TaxonomyCreate(ctx context.Context, taxonomy TaxonomyInterface) error

	// TaxonomyDelete permanently removes a taxonomy from the store.
	TaxonomyDelete(ctx context.Context, taxonomy TaxonomyInterface) error

	// TaxonomyFindByID retrieves a taxonomy by its unique identifier.
	TaxonomyFindByID(ctx context.Context, id string) (TaxonomyInterface, error)

	// TaxonomyFindBySlug retrieves a taxonomy by its URL-friendly slug.
	TaxonomyFindBySlug(ctx context.Context, slug string) (TaxonomyInterface, error)

	// TaxonomyList retrieves taxonomies matching the provided query options.
	TaxonomyList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyInterface, error)

	// TaxonomyUpdate modifies an existing taxonomy.
	TaxonomyUpdate(ctx context.Context, taxonomy TaxonomyInterface) error

	// Term methods manage individual terms within taxonomies.

	// TermCount returns the number of terms matching the query options.
	TermCount(ctx context.Context, options TermQueryOptions) (int64, error)

	// TermCreate inserts a new term into the store.
	TermCreate(ctx context.Context, term TermInterface) error

	// TermDelete permanently removes a term from the store.
	TermDelete(ctx context.Context, term TermInterface) error

	// TermFindByID retrieves a term by its unique identifier.
	TermFindByID(ctx context.Context, id string) (TermInterface, error)

	// TermFindBySlug retrieves a term by its taxonomy slug and term slug.
	TermFindBySlug(ctx context.Context, taxonomySlug, termSlug string) (TermInterface, error)

	// TermList retrieves terms matching the provided query options.
	TermList(ctx context.Context, options TermQueryOptions) ([]TermInterface, error)

	// TermUpdate modifies an existing term.
	TermUpdate(ctx context.Context, term TermInterface) error

	// Post-term relationship methods manage associations between posts and terms.

	// PostAddTerm appends a term to a post (adds at the end of the sequence).
	// Automatically calculates the next available sequence number.
	// Returns an error if taxonomy features are not enabled.
	PostAddTerm(ctx context.Context, postID string, termID string) error

	// PostInsertTermAt associates a term with a post at the specified sequence/order.
	PostInsertTermAt(ctx context.Context, postID string, termID string, sequence int) error

	// PostMoveTermTo moves a term to a specific sequence position on a post.
	// Reorders existing terms by pushing subsequent terms down (incrementing their sequence).
	// Returns an error if the term is not associated with the post.
	PostMoveTermTo(ctx context.Context, postID string, termID string, sequence int) error

	// PostRemoveTerm dissociates a term from a post.
	PostRemoveTerm(ctx context.Context, postID string, termID string) error

	// TermListByPostID retrieves all terms associated with a post for a specific taxonomy.
	TermListByPostID(ctx context.Context, postID string, taxonomySlug string) ([]TermInterface, error)

	// PostSetTerms replaces all terms for a post within a taxonomy with the provided term IDs.
	PostSetTerms(ctx context.Context, postID string, taxonomySlug string, termIDs []string) error

	// PostListByTermID retrieves all posts associated with a specific term ID.
	PostListByTermID(ctx context.Context, termID string, options PostQueryOptions) ([]PostInterface, error)

	// Utility methods provide helper operations for term management.

	// TermIncrementCount increases the usage count for a term.
	TermIncrementCount(ctx context.Context, termID string) error

	// TermDecrementCount decreases the usage count for a term.
	TermDecrementCount(ctx context.Context, termID string) error
}

var _ StoreInterface = (*storeImplementation)(nil) // verify it extends the interface

// storeImplementation is the concrete implementation of the StoreInterface.
// It provides database operations for posts, taxonomies, terms, and term relations.
type storeImplementation struct {
	postTableName         string
	taxonomyTableName     string
	termTableName         string
	termRelationTableName string
	db                    *neat.Database
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
	// Use raw SQL to add column if it doesn't exist
	// This is a temporary migration for existing installations
	sql := `ALTER TABLE ` + store.postTableName + ` ADD COLUMN ` + COLUMN_SLUG + ` VARCHAR(255)`

	// Get underlying DB to execute raw SQL
	db, err := store.db.DB()
	if err != nil {
		return err
	}

	// Try to execute, ignore error if column already exists
	_, err = db.Exec(sql)
	if err != nil {
		// Column might already exist, which is fine
		return nil
	}
	return nil
}

// MigrateUp creates the blog store tables
func (store *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	// Create main post table
	if !store.db.Schema().HasTable(store.postTableName) {
		err := store.db.Schema().Create(store.postTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 21)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_SLUG, 255)
			table.Text(COLUMN_TITLE)
			table.Text(COLUMN_CONTENT)
			table.Text(COLUMN_SUMMARY)
			table.String(COLUMN_STATUS, 50)
			table.String(COLUMN_AUTHOR_ID, 40)
			table.String(COLUMN_CANONICAL_URL, 255)
			table.String(COLUMN_IMAGE_URL, 255)
			table.String(COLUMN_MEMO, 255)
			table.String(COLUMN_META_DESCRIPTION, 255)
			table.String(COLUMN_META_KEYWORDS, 255)
			table.String(COLUMN_META_ROBOTS, 50)
			table.Text(COLUMN_METAS)
			table.Boolean(COLUMN_FEATURED)
			table.DateTime(COLUMN_PUBLISHED_AT)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(constants.SoftDeleteAtColumn).Default(constants.MaxSoftDeletedAtDefault)
		})
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
	}

	// Create taxonomy tables only if enabled
	if store.taxonomyEnabled {
		// Create taxonomy table
		if !store.db.Schema().HasTable(store.taxonomyTableName) {
			err := store.db.Schema().Create(store.taxonomyTableName, func(table contractsschema.Blueprint) {
				table.String(COLUMN_ID, 21)
				table.Primary(COLUMN_ID)
				table.String(COLUMN_NAME, 255)
				table.String(COLUMN_SLUG, 255)
				table.Text(COLUMN_DESCRIPTION)
				table.DateTime(COLUMN_CREATED_AT)
				table.DateTime(COLUMN_UPDATED_AT)
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}

		// Create term table
		if !store.db.Schema().HasTable(store.termTableName) {
			err := store.db.Schema().Create(store.termTableName, func(table contractsschema.Blueprint) {
				table.String(COLUMN_ID, 21)
				table.Primary(COLUMN_ID)
				table.String(COLUMN_TAXONOMY_ID, 21)
				table.String(COLUMN_PARENT_ID, 21)
				table.Integer(COLUMN_SEQUENCE)
				table.String(COLUMN_NAME, 255)
				table.String(COLUMN_SLUG, 255)
				table.Text(COLUMN_DESCRIPTION)
				table.Integer(COLUMN_COUNT)
				table.DateTime(COLUMN_CREATED_AT)
				table.DateTime(COLUMN_UPDATED_AT)
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}

		// Create term relation table
		if !store.db.Schema().HasTable(store.termRelationTableName) {
			err := store.db.Schema().Create(store.termRelationTableName, func(table contractsschema.Blueprint) {
				table.String(COLUMN_ID, 21)
				table.Primary(COLUMN_ID)
				table.String(COLUMN_POST_ID, 21)
				table.String(COLUMN_TERM_ID, 21)
				table.Integer(COLUMN_SEQUENCE)
				table.DateTime(COLUMN_CREATED_AT)
				table.DateTime(COLUMN_UPDATED_AT)
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	if store.versioningEnabled {
		if store.versioningStore == nil {
			return errors.New("versioning store is nil")
		}
		if err := store.versioningStore.MigrateUp(ctx, tx...); err != nil {
			return err
		}
	}

	return nil
}

// MigrateDown drops the blog store tables
func (store *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	// Drop tables in reverse order of creation (due to potential foreign key constraints)
	if store.taxonomyEnabled {
		// Drop term relation table first
		if store.db.Schema().HasTable(store.termRelationTableName) {
			err := store.db.Schema().Drop(store.termRelationTableName)
			if err != nil {
				log.Println(err)
				return err
			}
		}

		// Drop term table
		if store.db.Schema().HasTable(store.termTableName) {
			err := store.db.Schema().Drop(store.termTableName)
			if err != nil {
				log.Println(err)
				return err
			}
		}

		// Drop taxonomy table
		if store.db.Schema().HasTable(store.taxonomyTableName) {
			err := store.db.Schema().Drop(store.taxonomyTableName)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	// Drop post table
	if store.db.Schema().HasTable(store.postTableName) {
		err := store.db.Schema().Drop(store.postTableName)
		if err != nil {
			log.Println(err)
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

// GetPostTableName returns the post table name
func (st *storeImplementation) GetPostTableName() string {
	return st.postTableName
}

// SetPostTableName sets the post table name
func (st *storeImplementation) SetPostTableName(tableName string) {
	st.postTableName = tableName
}

// GetTaxonomyTableName returns the taxonomy table name
func (st *storeImplementation) GetTaxonomyTableName() string {
	return st.taxonomyTableName
}

// SetTaxonomyTableName sets the taxonomy table name
func (st *storeImplementation) SetTaxonomyTableName(tableName string) {
	st.taxonomyTableName = tableName
}

// GetTermTableName returns the term table name
func (st *storeImplementation) GetTermTableName() string {
	return st.termTableName
}

// SetTermTableName sets the term table name
func (st *storeImplementation) SetTermTableName(tableName string) {
	st.termTableName = tableName
}

// GetTermRelationTableName returns the term relation table name
func (st *storeImplementation) GetTermRelationTableName() string {
	return st.termRelationTableName
}

// SetTermRelationTableName sets the term relation table name
func (st *storeImplementation) SetTermRelationTableName(tableName string) {
	st.termRelationTableName = tableName
}

// PostCreate inserts a new post into the database.
// It sets the created_at and updated_at timestamps automatically.
// Also tracks the creation in the versioning store if versioning is enabled.
func (store *storeImplementation) PostCreate(ctx context.Context, post PostInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if post.GetID() == "" {
		post.SetID(GenerateShortID())
	}

	post.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	post.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	db, err := store.db.DB()
	if err != nil {
		return err
	}

	metas, _ := post.GetMetas()
	metasJSON := ""
	if len(metas) > 0 {
		metasBytes, err := json.Marshal(metas)
		if err != nil {
			return err
		}
		metasJSON = string(metasBytes)
	}

	_, err = db.ExecContext(ctx, "INSERT INTO "+store.postTableName+" (id, slug, title, content, summary, status, author_id, canonical_url, image_url, memo, meta_description, meta_keywords, meta_robots, metas, featured, published_at, created_at, updated_at, soft_deleted_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		post.GetID(),
		post.GetSlug(),
		post.GetTitle(),
		post.GetContent(),
		post.GetSummary(),
		post.GetStatus(),
		post.GetAuthorID(),
		post.GetCanonicalURL(),
		post.GetImageUrl(),
		post.GetMemo(),
		post.GetMetaDescription(),
		post.GetMetaKeywords(),
		post.GetMetaRobots(),
		metasJSON,
		post.GetFeatured(),
		post.GetPublishedAtCarbon().StdTime(),
		post.GetCreatedAtCarbon().StdTime(),
		post.GetUpdatedAtCarbon().StdTime(),
		post.GetSoftDeletedAtCarbon().StdTime(),
	)

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
	if ctx == nil {
		return 0, errors.New("ctx is nil")
	}

	q := store.buildPostQuery(options)

	var count int64
	err := q.Table(store.postTableName).Count(&count)
	return count, err
}

// PostTrash moves a post to trash by setting its status to POST_STATUS_TRASH.
func (store *storeImplementation) PostTrash(ctx context.Context, post PostInterface) error {
	post.SetStatus(POST_STATUS_TRASH)

	return store.PostUpdate(ctx, post)
}

// PostDelete permanently removes a post from the database.
func (store *storeImplementation) PostDelete(ctx context.Context, post PostInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if post == nil {
		return errors.New("post is nil")
	}

	return store.PostDeleteByID(ctx, post.GetID())
}

// PostDeleteByID permanently removes a post by its ID.
func (store *storeImplementation) PostDeleteByID(ctx context.Context, id string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("post id is empty")
	}

	_, err := store.db.Query().
		Table(store.postTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()

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

// PostFindBySlug retrieves a post by its slug.
func (store *storeImplementation) PostFindBySlug(ctx context.Context, slug string) (PostInterface, error) {
	if slug == "" {
		return nil, errors.New("slug is empty")
	}

	list, err := store.PostList(ctx, PostQueryOptions{
		Slug:  slug,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// PostFindByOldSlug retrieves a post by its old slug (for redirect handling).
func (store *storeImplementation) PostFindByOldSlug(ctx context.Context, oldSlug string) (PostInterface, error) {
	if oldSlug == "" {
		return nil, errors.New("old slug is empty")
	}

	list, err := store.PostList(ctx, PostQueryOptions{
		OldSlug: oldSlug,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// PostFindPrevious finds the post created immediately before the given post.
func (st *storeImplementation) PostFindPrevious(post PostInterface) (PostInterface, error) {
	list, err := st.PostList(context.Background(), PostQueryOptions{
		CreatedAtLessThan: post.GetCreatedAtCarbon().ToDateTimeString(),
		Limit:             1,
		OrderBy:           COLUMN_CREATED_AT,
		SortOrder:         "DESC",
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
		OrderBy:              COLUMN_CREATED_AT,
		SortOrder:            "ASC",
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
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	type postRow struct {
		ID              string    `db:"id"`
		Slug            string    `db:"slug"`
		Title           string    `db:"title"`
		Content         string    `db:"content"`
		Summary         string    `db:"summary"`
		Status          string    `db:"status"`
		AuthorID        string    `db:"author_id"`
		CanonicalURL    string    `db:"canonical_url"`
		ImageURL        string    `db:"image_url"`
		Memo            string    `db:"memo"`
		MetaDescription string    `db:"meta_description"`
		MetaKeywords    string    `db:"meta_keywords"`
		MetaRobots      string    `db:"meta_robots"`
		Metas           string    `db:"metas"`
		Featured        string    `db:"featured"`
		PublishedAt     time.Time `db:"published_at"`
		CreatedAt       time.Time `db:"created_at"`
		UpdatedAt       time.Time `db:"updated_at"`
		SoftDeletedAt   time.Time `db:"soft_deleted_at"`
	}

	q := st.buildPostQuery(options)

	var rows []postRow
	if err := q.Table(st.postTableName).Get(&rows); err != nil {
		return []PostInterface{}, err
	}

	list := make([]PostInterface, 0, len(rows))
	for _, r := range rows {
		p := NewPost()
		p.SetID(r.ID)
		p.SetSlug(r.Slug)
		p.SetTitle(r.Title)
		p.SetContent(r.Content)
		p.SetSummary(r.Summary)
		p.SetStatus(r.Status)
		p.SetAuthorID(r.AuthorID)
		p.SetCanonicalURL(r.CanonicalURL)
		p.SetImageUrl(r.ImageURL)
		p.SetMemo(r.Memo)
		p.SetMetaDescription(r.MetaDescription)
		p.SetMetaKeywords(r.MetaKeywords)
		p.SetMetaRobots(r.MetaRobots)
		// Parse JSON string to map for SetMetas
		if r.Metas != "" {
			var metas map[string]string
			if err := json.Unmarshal([]byte(r.Metas), &metas); err == nil {
				for k, v := range metas {
					p.SetMeta(k, v)
				}
			}
		}
		p.SetFeatured(r.Featured)
		if postImpl, ok := p.(*postImplementation); ok {
			postImpl.PublishedAtField = r.PublishedAt
			postImpl.CreatedAtField.CreatedAt = r.CreatedAt
			postImpl.UpdatedAtField.UpdatedAt = r.UpdatedAt
			postImpl.SoftDeletedAt = r.SoftDeletedAt
		}
		list = append(list, p)
	}

	return list, nil
}

// PostSoftDelete marks a post as deleted by setting the soft_deleted_at timestamp.
func (st *storeImplementation) PostSoftDelete(ctx context.Context, post PostInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if post == nil {
		return errors.New("post is nil")
	}

	post.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return st.PostUpdate(ctx, post)
}

// PostSoftDeleteByID marks a post as deleted by its ID.
func (st *storeImplementation) PostSoftDeleteByID(ctx context.Context, id string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	post, err := st.PostFindByID(ctx, id)

	if err != nil {
		return err
	}

	if post == nil {
		return errors.New("post not found")
	}

	return st.PostSoftDelete(ctx, post)
}

// PostUpdate updates an existing post in the database.
// Only changed fields are updated. Also tracks the update in the versioning store if enabled.
func (st *storeImplementation) PostUpdate(ctx context.Context, post PostInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if post == nil {
		return errors.New("post is nil")
	}

	post.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := post.GetDataChanged()

	delete(dataChanged, "id")   // ID is not updatable
	delete(dataChanged, "hash") // Hash is not updatable
	delete(dataChanged, "data") // Data is not updatable

	if len(dataChanged) < 1 {
		return nil
	}

	// Convert dataChanged to proper format for neat Update
	// neat expects map[string]interface{} with proper Go types
	updateData := make(map[string]interface{})
	for k, v := range dataChanged {
		updateData[k] = v
	}

	// Handle special fields that need conversion
	if publishedAt, ok := updateData["published_at"]; ok {
		if publishedAtStr, ok := publishedAt.(string); ok {
			updateData["published_at"] = carbon.Parse(publishedAtStr, carbon.UTC).StdTime()
		}
	}
	if createdAt, ok := updateData["created_at"]; ok {
		if createdAtStr, ok := createdAt.(string); ok {
			updateData["created_at"] = carbon.Parse(createdAtStr, carbon.UTC).StdTime()
		}
	}
	if updatedAt, ok := updateData["updated_at"]; ok {
		if updatedAtStr, ok := updatedAt.(string); ok {
			updateData["updated_at"] = carbon.Parse(updatedAtStr, carbon.UTC).StdTime()
		}
	}
	if softDeletedAt, ok := updateData["soft_deleted_at"]; ok {
		if softDeletedAtStr, ok := softDeletedAt.(string); ok {
			updateData["soft_deleted_at"] = carbon.Parse(softDeletedAtStr, carbon.UTC).StdTime()
		}
	}

	_, err := st.db.Query().
		Table(st.postTableName).
		Where(COLUMN_ID+" = ?", post.GetID()).
		Update(updateData)

	if err != nil {
		return err
	}

	post.MarkAsNotDirty()
	if err2 := st.versioningTrackEntity(ctx, VERSIONING_TYPE_POST, post.GetID(), post); err2 != nil {
		return err2
	}

	return nil
}

// buildPostQuery builds a neat query from the post query options.
func (st *storeImplementation) buildPostQuery(options PostQueryOptions) contractsorm.Query {
	q := st.db.Query()

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if len(options.IDIn) > 0 {
		// Build IN clause manually for neat compatibility
		inClause := COLUMN_ID + " IN ("
		placeholders := make([]interface{}, 0, len(options.IDIn))
		for i, id := range options.IDIn {
			if i > 0 {
				inClause += ", "
			}
			inClause += "?"
			placeholders = append(placeholders, id)
		}
		inClause += ")"
		q = q.Where(inClause, placeholders...)
	}

	if options.Slug != "" {
		q = q.Where(COLUMN_SLUG+" = ?", options.Slug)
	}

	if options.OldSlug != "" {
		// Use JSON contains to check if old slugs array contains the value
		// The JSON structure is: {"_old_slugs": "[\"slug1\",\"slug2\"]"}
		// Note: the value is a string containing JSON array
		// We search for the pattern: "_old_slugs":"[..."old-slug-1"...]"
		// Use escaped quotes for the pattern
		q = q.Where(COLUMN_METAS+" LIKE ?", "%\"_old_slugs\":\"[%\\\""+options.OldSlug+"\\\"%]%")
	}

	if len(options.MetaEquals) > 0 {
		// For each meta key-value pair, add a JSON contains condition
		// The JSON structure is: {"key": "value"}
		for key, value := range options.MetaEquals {
			// Search for pattern: "key":"value"
			q = q.Where(COLUMN_METAS+" LIKE ?", "%\""+key+"\":\""+value+"\"%")
		}
	}

	if len(options.MetaArrayContains) > 0 {
		// For each meta array key-value pair, add a JSON contains condition
		// The JSON structure is: {"key": "[\"value1\",\"value2\"]"}
		for key, value := range options.MetaArrayContains {
			// Search for pattern: "key":"[..."value"...]"
			q = q.Where(COLUMN_METAS+" LIKE ?", "%\""+key+"\":\"[%\""+value+"\"%]%")
		}
	}

	if options.Status != "" {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status)
	}

	if len(options.StatusIn) > 0 {
		// Build IN clause manually for neat compatibility
		inClause := COLUMN_STATUS + " IN ("
		placeholders := make([]interface{}, 0, len(options.StatusIn))
		for i, status := range options.StatusIn {
			if i > 0 {
				inClause += ", "
			}
			inClause += "?"
			placeholders = append(placeholders, status)
		}
		inClause += ")"
		q = q.Where(inClause, placeholders...)
	}

	if options.CreatedAtLessThan != "" {
		q = q.Where(COLUMN_CREATED_AT+" < ?", carbon.Parse(options.CreatedAtLessThan, carbon.UTC).StdTime())
	}

	if options.CreatedAtGreaterThan != "" {
		q = q.Where(COLUMN_CREATED_AT+" > ?", carbon.Parse(options.CreatedAtGreaterThan, carbon.UTC).StdTime())
	}

	if options.Search != "" {
		// Simple search on title and content
		q = q.Where("("+COLUMN_TITLE+" LIKE ? OR "+COLUMN_CONTENT+" LIKE ?)", "%"+options.Search+"%", "%"+options.Search+"%")
	}

	if options.OrderBy != "" {
		order := options.SortOrder
		if order == "" {
			order = "DESC"
		}
		q = q.OrderBy(options.OrderBy + " " + order)
	}

	if options.Limit > 0 {
		q = q.Limit(options.Limit)
	}

	if options.Offset > 0 {
		q = q.Offset(options.Offset)
	}

	// Handle soft delete filtering
	if options.WithDeleted {
		q = q.WithSoftDeleted()
	} else {
		// By default, filter out soft-deleted records
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", carbon.Parse(MAX_DATETIME, carbon.UTC).StdTime())
	}

	return q
}
