package blogstore

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// VersioningInterface represents a version entry for tracking entity changes over time.
type VersioningInterface interface {
	IsSoftDeleted() bool

	ID() string
	SetID(id string) VersioningInterface

	EntityType() string
	SetEntityType(entityType string) VersioningInterface

	EntityID() string
	SetEntityID(entityID string) VersioningInterface

	Content() string
	SetContent(content string) VersioningInterface

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) VersioningInterface

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(softDeletedAt string) VersioningInterface
}

// VersioningQueryInterface provides query options for retrieving version entries.
type VersioningQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) VersioningQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) VersioningQueryInterface

	HasID() bool
	ID() string
	SetID(id string) VersioningQueryInterface

	HasEntityID() bool
	EntityID() string
	SetEntityID(entityID string) VersioningQueryInterface

	HasEntityType() bool
	EntityType() string
	SetEntityType(entityType string) VersioningQueryInterface

	HasOffset() bool
	Offset() int64
	SetOffset(offset int64) VersioningQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) VersioningQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) VersioningQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) VersioningQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(includeSoftDeleted bool) VersioningQueryInterface
}

// NewVersioning creates a new VersioningInterface instance.
// This is used to create a new version entry for tracking entity changes.
func NewVersioning() VersioningInterface {
	return NewVersioningFromExistingData(nil)
}

// NewVersioningFromExistingData creates a versioning from existing data.
func NewVersioningFromExistingData(data map[string]string) VersioningInterface {
	o := &versioningImplementation{}

	if data == nil {
		o.SetID(GenerateShortID())
		o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
		o.SetSoftDeletedAt(MAX_DATETIME)
		return o
	}

	o.SetID(data[COLUMN_ID])
	o.SetEntityType(data[COLUMN_ENTITY_TYPE])
	o.SetEntityID(data[COLUMN_ENTITY_ID])
	o.SetContent(data[COLUMN_CONTENT])
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		o.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		o.SetSoftDeletedAt(v)
	}
	return o
}

// versioningImplementation implements VersioningInterface.
type versioningImplementation struct {
	IDField         string
	EntityTypeField string
	EntityIDField   string
	ContentField    string
	CreatedAt       time.Time
	SoftDeletedAt   time.Time
}

var _ VersioningInterface = (*versioningImplementation)(nil)

// IsSoftDeleted returns true if the version is soft deleted.
func (o *versioningImplementation) IsSoftDeleted() bool {
	return o.SoftDeletedAt.Before(time.Now().UTC())
}

// ID returns the id of the version.
func (o *versioningImplementation) ID() string {
	return o.IDField
}

// SetID sets the id of the version.
func (o *versioningImplementation) SetID(id string) VersioningInterface {
	o.IDField = id
	return o
}

// EntityType returns the entity type of the version.
func (o *versioningImplementation) EntityType() string {
	return o.EntityTypeField
}

// SetEntityType sets the entity type of the version.
func (o *versioningImplementation) SetEntityType(entityType string) VersioningInterface {
	o.EntityTypeField = entityType
	return o
}

// EntityID returns the entity id of the version.
func (o *versioningImplementation) EntityID() string {
	return o.EntityIDField
}

// SetEntityID sets the entity id of the version.
func (o *versioningImplementation) SetEntityID(entityID string) VersioningInterface {
	o.EntityIDField = entityID
	return o
}

// Content returns the content of the version.
func (o *versioningImplementation) Content() string {
	return o.ContentField
}

// SetContent sets the content of the version.
func (o *versioningImplementation) SetContent(content string) VersioningInterface {
	o.ContentField = content
	return o
}

// GetCreatedAt returns the created at time of the version.
func (o *versioningImplementation) GetCreatedAt() string {
	if o.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAt).ToDateTimeString()
}

// GetCreatedAtCarbon returns the created at time of the version as a carbon object.
func (o *versioningImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAt)
}

// SetCreatedAt sets the created at time of the version.
func (o *versioningImplementation) SetCreatedAt(createdAt string) VersioningInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// GetSoftDeletedAt returns the soft deleted at time of the version.
func (o *versioningImplementation) GetSoftDeletedAt() string {
	if o.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletedAt).ToDateTimeString()
}

// GetSoftDeletedAtCarbon returns the soft deleted at time of the version as a carbon object.
func (o *versioningImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletedAt)
}

// SetSoftDeletedAt sets the soft deleted at time of the version.
func (o *versioningImplementation) SetSoftDeletedAt(softDeletedAt string) VersioningInterface {
	if softDeletedAt == "" {
		return o
	}
	o.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	return o
}
