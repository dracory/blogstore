package blogstore

import (
	"database/sql"
	"errors"

	"github.com/dracory/sb"
)

// NewStoreOptions define the options for creating a new block store
type NewStoreOptions struct {
	PostTableName      string
	DB                 *sql.DB
	DbDriverName       string
	TimeoutSeconds     int64
	AutomigrateEnabled bool
	DebugEnabled       bool
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

	store := &store{
		postTableName:      opts.PostTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
	}

	store.timeoutSeconds = 2 * 60 * 60 // 2 hours

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}
