package blogstore

import "github.com/dracory/versionstore"

// VersioningInterface is an alias for versionstore.VersionInterface.
// It represents a version entry for tracking entity changes over time.
type VersioningInterface = versionstore.VersionInterface

// VersioningQueryInterface is an alias for versionstore.VersionQueryInterface.
// It provides query options for retrieving version entries.
type VersioningQueryInterface = versionstore.VersionQueryInterface

// NewVersioning creates a new VersioningInterface instance.
// This is used to create a new version entry for tracking entity changes.
func NewVersioning() VersioningInterface {
	return versionstore.NewVersion()
}

// NewVersioningQuery creates a new VersioningQueryInterface instance.
// This is used to query version entries with filtering and sorting options.
func NewVersioningQuery() VersioningQueryInterface {
	return versionstore.NewVersionQuery()
}
