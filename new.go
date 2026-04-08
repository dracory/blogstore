package blogstore

import (
	"database/sql"
	"errors"

	"github.com/dracory/sb"
	"github.com/dracory/versionstore"
)

// NewStoreOptions defines the configuration options for creating a new blog store.
type NewStoreOptions struct {
	PostTableName         string
	TaxonomyTableName     string
	TermTableName         string
	TermRelationTableName string
	DB                    *sql.DB
	DbDriverName          string
	TimeoutSeconds        int64
	AutomigrateEnabled    bool
	DebugEnabled          bool

	VersioningEnabled   bool
	VersioningTableName string

	TaxonomyEnabled bool
}

// NewStore creates a new blog store with the provided options.
// It validates required fields, sets defaults for optional fields, and optionally runs AutoMigrate.
// Returns a StoreInterface for interacting with posts, taxonomies, and terms.
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.PostTableName == "" {
		return nil, errors.New("blog store: PostTableName is required")
	}

	if opts.TaxonomyTableName == "" {
		opts.TaxonomyTableName = "blog_taxonomy"
	}

	if opts.TermTableName == "" {
		opts.TermTableName = "blog_term"
	}

	if opts.TermRelationTableName == "" {
		opts.TermRelationTableName = "blog_term_rel"
	}

	if opts.DB == nil {
		return nil, errors.New("blog store: DB is required")
	}

	if opts.DbDriverName == "" {
		opts.DbDriverName = sb.DatabaseDriverName(opts.DB)
	}

	if opts.VersioningEnabled && opts.VersioningTableName == "" {
		return nil, errors.New("blog store: VersioningTableName is required")
	}

	var versionStore versionstore.StoreInterface
	if opts.VersioningEnabled {
		vs, err := versionstore.NewStore(versionstore.NewStoreOptions{
			TableName:          opts.VersioningTableName,
			DB:                 opts.DB,
			AutomigrateEnabled: opts.AutomigrateEnabled,
			DebugEnabled:       opts.DebugEnabled,
		})
		if err != nil {
			return nil, err
		}
		if vs == nil {
			return nil, errors.New("blog store: version store is nil")
		}
		versionStore = vs
	}

	store := &storeImplementation{
		postTableName:         opts.PostTableName,
		taxonomyTableName:     opts.TaxonomyTableName,
		termTableName:         opts.TermTableName,
		termRelationTableName: opts.TermRelationTableName,
		automigrateEnabled:    opts.AutomigrateEnabled,
		db:                    opts.DB,
		dbDriverName:          opts.DbDriverName,
		debugEnabled:          opts.DebugEnabled,
		versioningEnabled:     opts.VersioningEnabled,
		versioningStore:       versionStore,
		taxonomyEnabled:       opts.TaxonomyEnabled,
	}

	store.timeoutSeconds = 2 * 60 * 60 // 2 hours

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}
