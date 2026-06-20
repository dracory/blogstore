package blogstore

import "errors"

// NewVersioningQuery creates a new VersioningQueryInterface instance.
// This is used to query version entries with filtering and sorting options.
func NewVersioningQuery() VersioningQueryInterface {
	return &versioningQueryImplementation{
		properties: map[string]any{},
	}
}

// versioningQueryImplementation implements VersioningQueryInterface.
type versioningQueryImplementation struct {
	properties map[string]any
}

var _ VersioningQueryInterface = (*versioningQueryImplementation)(nil)

// Validate validates the query parameters.
func (q *versioningQueryImplementation) Validate() error {
	if q.HasEntityID() && q.EntityID() == "" {
		return errors.New("version query. entity_id cannot be empty")
	}

	if q.HasEntityType() && q.EntityType() == "" {
		return errors.New("version query. entity_type cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("version query. id cannot be empty")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("version query. limit cannot be negative")
	}

	if q.HasLimit() && q.Limit() < 1 {
		return errors.New("version query. limit cannot be less than 1")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("version query. offset cannot be negative")
	}

	return nil
}

// Columns returns the columns to select.
func (q *versioningQueryImplementation) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

// SetColumns sets the columns to select.
func (q *versioningQueryImplementation) SetColumns(columns []string) VersioningQueryInterface {
	q.properties["columns"] = columns
	return q
}

// IsCountOnly returns true if only count is requested.
func (q *versioningQueryImplementation) IsCountOnly() bool {
	if !q.hasProperty("count_only") {
		return false
	}

	return q.properties["count_only"].(bool)
}

// HasCountOnly returns true if count_only is set.
func (q *versioningQueryImplementation) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

// SetCountOnly sets whether to return only count.
func (q *versioningQueryImplementation) SetCountOnly(countOnly bool) VersioningQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

// HasEntityID returns true if entity_id is set.
func (q *versioningQueryImplementation) HasEntityID() bool {
	return q.hasProperty("entity_id")
}

// EntityID returns the entity ID.
func (q *versioningQueryImplementation) EntityID() string {
	if !q.hasProperty("entity_id") {
		return ""
	}

	return q.properties["entity_id"].(string)
}

// SetEntityID sets the entity ID.
func (q *versioningQueryImplementation) SetEntityID(entityID string) VersioningQueryInterface {
	q.properties["entity_id"] = entityID
	return q
}

// HasEntityType returns true if entity_type is set.
func (q *versioningQueryImplementation) HasEntityType() bool {
	return q.hasProperty("entity_type")
}

// EntityType returns the entity type.
func (q *versioningQueryImplementation) EntityType() string {
	if !q.hasProperty("entity_type") {
		return ""
	}

	return q.properties["entity_type"].(string)
}

// SetEntityType sets the entity type.
func (q *versioningQueryImplementation) SetEntityType(entityType string) VersioningQueryInterface {
	q.properties["entity_type"] = entityType
	return q
}

// HasID returns true if id is set.
func (q *versioningQueryImplementation) HasID() bool {
	return q.hasProperty("id")
}

// ID returns the version ID.
func (q *versioningQueryImplementation) ID() string {
	if !q.hasProperty("id") {
		return ""
	}

	return q.properties["id"].(string)
}

// SetID sets the version ID.
func (q *versioningQueryImplementation) SetID(id string) VersioningQueryInterface {
	q.properties["id"] = id
	return q
}

// HasLimit returns true if limit is set.
func (q *versioningQueryImplementation) HasLimit() bool {
	return q.hasProperty("limit")
}

// Limit returns the query limit.
func (q *versioningQueryImplementation) Limit() int {
	if !q.hasProperty("limit") {
		return 0
	}

	return q.properties["limit"].(int)
}

// SetLimit sets the query limit.
func (q *versioningQueryImplementation) SetLimit(limit int) VersioningQueryInterface {
	q.properties["limit"] = limit
	return q
}

// HasOffset returns true if offset is set.
func (q *versioningQueryImplementation) HasOffset() bool {
	return q.hasProperty("offset")
}

// Offset returns the query offset.
func (q *versioningQueryImplementation) Offset() int64 {
	if !q.hasProperty("offset") {
		return 0
	}

	return q.properties["offset"].(int64)
}

// SetOffset sets the query offset.
func (q *versioningQueryImplementation) SetOffset(offset int64) VersioningQueryInterface {
	q.properties["offset"] = offset
	return q
}

// HasOrderBy returns true if order_by is set.
func (q *versioningQueryImplementation) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

// OrderBy returns the order by field.
func (q *versioningQueryImplementation) OrderBy() string {
	if !q.hasProperty("order_by") {
		return ""
	}

	return q.properties["order_by"].(string)
}

// SetOrderBy sets the order by field.
func (q *versioningQueryImplementation) SetOrderBy(orderBy string) VersioningQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

// HasSortOrder returns true if sort_order is set.
func (q *versioningQueryImplementation) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

// SortOrder returns the sort order (ASC or DESC).
func (q *versioningQueryImplementation) SortOrder() string {
	if !q.hasProperty("sort_order") {
		return ""
	}

	return q.properties["sort_order"].(string)
}

// SetSortOrder sets the sort order (ASC or DESC).
func (q *versioningQueryImplementation) SetSortOrder(sortOrder string) VersioningQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

// HasSoftDeletedIncluded returns true if soft_deleted_included is set.
func (q *versioningQueryImplementation) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_deleted_included")
}

// SoftDeletedIncluded returns true if soft deleted versions should be included.
func (q *versioningQueryImplementation) SoftDeletedIncluded() bool {
	if q.hasProperty("soft_deleted_included") {
		return q.properties["soft_deleted_included"].(bool)
	}

	return false
}

// SetSoftDeletedIncluded sets whether to include soft deleted versions.
func (q *versioningQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) VersioningQueryInterface {
	q.properties["soft_deleted_included"] = softDeletedIncluded
	return q
}

// hasProperty returns true if the property exists in the map.
func (q *versioningQueryImplementation) hasProperty(key string) bool {
	return q.properties[key] != nil
}
