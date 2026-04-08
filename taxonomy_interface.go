package blogstore

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// TaxonomyInterface defines the interface for taxonomy types (category, tag, etc.).
// Taxonomies are used to group and classify content within the blog system.
type TaxonomyInterface interface {
	// Identity
	// GetID returns the unique identifier of the taxonomy.
	GetID() string
	// SetID sets the unique identifier of the taxonomy.
	SetID(id string) TaxonomyInterface

	// Display
	// GetName returns the display name of the taxonomy.
	GetName() string
	// SetName sets the display name of the taxonomy.
	SetName(name string) TaxonomyInterface

	// URL slug (e.g., "category", "tag")
	// GetSlug returns the URL-friendly slug of the taxonomy.
	GetSlug() string
	// SetSlug sets the URL-friendly slug of the taxonomy.
	SetSlug(slug string) TaxonomyInterface

	// Description
	// GetDescription returns the description of the taxonomy.
	GetDescription() string
	// SetDescription sets the description of the taxonomy.
	SetDescription(description string) TaxonomyInterface

	// Timestamps
	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) TaxonomyInterface
	// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
	GetCreatedAtTime() time.Time

	// GetUpdatedAt returns the last update timestamp as a string.
	GetUpdatedAt() string
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) TaxonomyInterface
	// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon
	// GetUpdatedAtTime returns the last update timestamp as a time.Time instance.
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

// TermInterface defines the interface for terms within taxonomies.
// Terms are individual items within a taxonomy (e.g., "Technology" within "Categories").
type TermInterface interface {
	// Identity
	// GetID returns the unique identifier of the term.
	GetID() string
	// SetID sets the unique identifier of the term.
	SetID(id string) TermInterface

	// Relationships
	// GetTaxonomyID returns the ID of the taxonomy this term belongs to.
	GetTaxonomyID() string
	// SetTaxonomyID sets the ID of the taxonomy this term belongs to.
	SetTaxonomyID(taxonomyID string) TermInterface

	// For hierarchy (empty if root)
	// GetParentID returns the ID of the parent term (for hierarchical terms).
	GetParentID() string
	// SetParentID sets the ID of the parent term (for hierarchical terms).
	SetParentID(parentID string) TermInterface

	// For ordering subcategories within parent (0 = default)
	// GetSequence returns the display sequence/order of the term.
	GetSequence() int
	// SetSequence sets the display sequence/order of the term.
	SetSequence(sequence int) TermInterface

	// Display
	// GetName returns the display name of the term.
	GetName() string
	// SetName sets the display name of the term.
	SetName(name string) TermInterface

	// URL slug
	// GetSlug returns the URL-friendly slug of the term.
	GetSlug() string
	// SetSlug sets the URL-friendly slug of the term.
	SetSlug(slug string) TermInterface

	// Description
	// GetDescription returns the description of the term.
	GetDescription() string
	// SetDescription sets the description of the term.
	SetDescription(description string) TermInterface

	// Cached post count
	// GetCount returns the number of posts associated with this term.
	GetCount() int
	// SetCount sets the number of posts associated with this term.
	SetCount(count int) TermInterface

	// Timestamps
	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) TermInterface
	// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
	GetCreatedAtTime() time.Time

	// GetUpdatedAt returns the last update timestamp as a string.
	GetUpdatedAt() string
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) TermInterface
	// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon
	// GetUpdatedAtTime returns the last update timestamp as a time.Time instance.
	GetUpdatedAtTime() time.Time

	// DataObject methods (from embedded dataobject.DataObject)
	// GetData returns all term data as a map.
	GetData() map[string]string
	// GetDataChanged returns only the fields that have been modified.
	GetDataChanged() map[string]string
	// MarkAsNotDirty clears the dirty state of the term.
	MarkAsNotDirty()
	// Get retrieves a value by key from the term data.
	Get(key string) string
	// Set stores a value by key in the term data.
	Set(key string, value string)
	// Hydrate populates the term with data from a map.
	Hydrate(data map[string]string)
	// IsDirty returns true if the term has unsaved changes.
	IsDirty() bool
}

// TermRelationInterface defines the interface for post-term relationships.
// Term relations link posts to taxonomy terms.
type TermRelationInterface interface {
	// Identity
	// GetID returns the unique identifier of the term relation.
	GetID() string
	// SetID sets the unique identifier of the term relation.
	SetID(id string) TermRelationInterface

	// Relationships
	// GetPostID returns the ID of the post in this relation.
	GetPostID() string
	// SetPostID sets the ID of the post in this relation.
	SetPostID(postID string) TermRelationInterface

	// GetTermID returns the ID of the term in this relation.
	GetTermID() string
	// SetTermID sets the ID of the term in this relation.
	SetTermID(termID string) TermRelationInterface

	// For manual ordering (0 = default)
	// GetSequence returns the display sequence/order of the relation.
	GetSequence() int
	// SetSequence sets the display sequence/order of the relation.
	SetSequence(sequence int) TermRelationInterface

	// Timestamps
	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) TermRelationInterface
	// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
	GetCreatedAtTime() time.Time

	// GetUpdatedAt returns the last update timestamp as a string.
	GetUpdatedAt() string
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) TermRelationInterface
	// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon
	// GetUpdatedAtTime returns the last update timestamp as a time.Time instance.
	GetUpdatedAtTime() time.Time

	// DataObject methods (from embedded dataobject.DataObject)
	// GetData returns all term relation data as a map.
	GetData() map[string]string
	// GetDataChanged returns only the fields that have been modified.
	GetDataChanged() map[string]string
	// MarkAsNotDirty clears the dirty state of the term relation.
	MarkAsNotDirty()
	// Get retrieves a value by key from the term relation data.
	Get(key string) string
	// Set stores a value by key in the term relation data.
	Set(key string, value string)
	// Hydrate populates the term relation with data from a map.
	Hydrate(data map[string]string)
	// IsDirty returns true if the term relation has unsaved changes.
	IsDirty() bool
}

// Compile-time checks to ensure implementations implement interfaces.
var _ TaxonomyInterface = (*taxonomyImplementation)(nil)
var _ TermInterface = (*termImplementation)(nil)
var _ TermRelationInterface = (*termRelationImplementation)(nil)
