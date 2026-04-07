package blogstore

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// TaxonomyInterface defines the interface for taxonomy types (category, tag, etc.)
type TaxonomyInterface interface {
	// Identity
	GetID() string
	SetID(id string) TaxonomyInterface

	// Display
	GetName() string
	SetName(name string) TaxonomyInterface

	// URL slug (e.g., "category", "tag")
	GetSlug() string
	SetSlug(slug string) TaxonomyInterface

	// Description
	GetDescription() string
	SetDescription(description string) TaxonomyInterface

	// Timestamps
	GetCreatedAt() string
	SetCreatedAt(createdAt string) TaxonomyInterface
	GetCreatedAtCarbon() *carbon.Carbon
	GetCreatedAtTime() time.Time

	GetUpdatedAt() string
	SetUpdatedAt(updatedAt string) TaxonomyInterface
	GetUpdatedAtCarbon() *carbon.Carbon
	GetUpdatedAtTime() time.Time

	// DataObject methods (from embedded dataobject.DataObject)
	GetData() map[string]string
	GetDataChanged() map[string]string
	MarkAsNotDirty()
	Get(key string) string
	Set(key string, value string)
	Hydrate(data map[string]string)
	IsDirty() bool
}

// TermInterface defines the interface for terms within taxonomies
type TermInterface interface {
	// Identity
	GetID() string
	SetID(id string) TermInterface

	// Relationships
	GetTaxonomyID() string
	SetTaxonomyID(taxonomyID string) TermInterface

	// For hierarchy (empty if root)
	GetParentID() string
	SetParentID(parentID string) TermInterface

	// For ordering subcategories within parent (0 = default)
	GetSequence() int
	SetSequence(sequence int) TermInterface

	// Display
	GetName() string
	SetName(name string) TermInterface

	// URL slug
	GetSlug() string
	SetSlug(slug string) TermInterface

	// Description
	GetDescription() string
	SetDescription(description string) TermInterface

	// Cached post count
	GetCount() int
	SetCount(count int) TermInterface

	// Timestamps
	GetCreatedAt() string
	SetCreatedAt(createdAt string) TermInterface
	GetCreatedAtCarbon() *carbon.Carbon
	GetCreatedAtTime() time.Time

	GetUpdatedAt() string
	SetUpdatedAt(updatedAt string) TermInterface
	GetUpdatedAtCarbon() *carbon.Carbon
	GetUpdatedAtTime() time.Time

	// DataObject methods (from embedded dataobject.DataObject)
	GetData() map[string]string
	GetDataChanged() map[string]string
	MarkAsNotDirty()
	Get(key string) string
	Set(key string, value string)
	Hydrate(data map[string]string)
	IsDirty() bool
}

// TermRelationInterface defines the interface for post-term relationships
type TermRelationInterface interface {
	// Identity
	GetID() string
	SetID(id string) TermRelationInterface

	// Relationships
	GetPostID() string
	SetPostID(postID string) TermRelationInterface

	GetTermID() string
	SetTermID(termID string) TermRelationInterface

	// For manual ordering (0 = default)
	GetSequence() int
	SetSequence(sequence int) TermRelationInterface

	// Timestamps
	GetCreatedAt() string
	SetCreatedAt(createdAt string) TermRelationInterface
	GetCreatedAtCarbon() *carbon.Carbon
	GetCreatedAtTime() time.Time

	GetUpdatedAt() string
	SetUpdatedAt(updatedAt string) TermRelationInterface
	GetUpdatedAtCarbon() *carbon.Carbon
	GetUpdatedAtTime() time.Time

	// DataObject methods (from embedded dataobject.DataObject)
	GetData() map[string]string
	GetDataChanged() map[string]string
	MarkAsNotDirty()
	Get(key string) string
	Set(key string, value string)
	Hydrate(data map[string]string)
	IsDirty() bool
}

// Ensure implementations implement interfaces
var _ TaxonomyInterface = (*taxonomyImplementation)(nil)
var _ TermInterface = (*termImplementation)(nil)
var _ TermRelationInterface = (*termRelationImplementation)(nil)
