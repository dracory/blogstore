package blogstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/neat"
)

// NewStoreOptions defines the configuration options for creating a new blog store.
type NewStoreOptions struct {
	PostTableName         string
	TaxonomyTableName     string
	TermTableName         string
	TermRelationTableName string
	PostFileTableName     string
	DB                    *sql.DB
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

	if opts.PostFileTableName == "" {
		opts.PostFileTableName = "blog_post_file"
	}

	if opts.DB == nil {
		return nil, errors.New("blog store: DB is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	if opts.VersioningEnabled && opts.VersioningTableName == "" {
		return nil, errors.New("blog store: VersioningTableName is required")
	}

	store := &storeImplementation{
		postTableName:         opts.PostTableName,
		taxonomyTableName:     opts.TaxonomyTableName,
		termTableName:         opts.TermTableName,
		termRelationTableName: opts.TermRelationTableName,
		postFileTableName:     opts.PostFileTableName,
		automigrateEnabled:    opts.AutomigrateEnabled,
		db:                    neatDB,
		debugEnabled:          opts.DebugEnabled,
		versioningEnabled:     opts.VersioningEnabled,
		versioningTableName:   opts.VersioningTableName,
		taxonomyEnabled:       opts.TaxonomyEnabled,
	}

	store.timeoutSeconds = 2 * 60 * 60 // 2 hours

	if store.automigrateEnabled {
		if err := store.MigrateUp(context.Background()); err != nil {
			return nil, err
		}
	}

	return store, nil
}
