package blogstore

import (
	"database/sql"
	"errors"

	"github.com/dracory/sb"
	"github.com/dracory/versionstore"
)

// NewStoreOptions define the options for creating a new block store
type NewStoreOptions struct {
	PostTableName      string
	DB                 *sql.DB
	DbDriverName       string
	TimeoutSeconds     int64
	AutomigrateEnabled bool
	DebugEnabled       bool

	VersioningEnabled   bool
	VersioningTableName string
}

// NewStore creates a new block store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.PostTableName == "" {
		return nil, errors.New("blog store: PostTableName is required")
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

	store := &store{
		postTableName:      opts.PostTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
		versioningEnabled:  opts.VersioningEnabled,
		versioningStore:    versionStore,
	}

	store.timeoutSeconds = 2 * 60 * 60 // 2 hours

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}
