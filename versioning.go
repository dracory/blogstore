package blogstore

import "github.com/dracory/versionstore"

type VersioningInterface = versionstore.VersionInterface

type VersioningQueryInterface = versionstore.VersionQueryInterface

func NewVersioning() VersioningInterface {
	return versionstore.NewVersion()
}

func NewVersioningQuery() VersioningQueryInterface {
	return versionstore.NewVersionQuery()
}
