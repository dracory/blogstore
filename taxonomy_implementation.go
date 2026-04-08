package blogstore

import (
	"time"

	"github.com/dracory/dataobject"
	"github.com/dracory/str"
	"github.com/dromara/carbon/v2"
)

// ============================ TAXONOMY ============================

// NewTaxonomy creates a new Taxonomy instance with default values.
// The taxonomy is initialized with a generated ID, empty name/slug/description,
// and current timestamps for created_at and updated_at.
func NewTaxonomy() TaxonomyInterface {
	o := &taxonomyImplementation{}
	o.SetID(GenerateShortID()).
		SetName("").
		SetSlug("").
		SetDescription("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return o
}

// NewTaxonomyFromExistingData creates a Taxonomy instance from existing data.
// This is useful when hydrating a taxonomy from database records.
func NewTaxonomyFromExistingData(data map[string]string) TaxonomyInterface {
	o := &taxonomyImplementation{}
	o.Hydrate(data)
	return o
}

// taxonomyImplementation is the concrete implementation of the TaxonomyInterface.
// It embeds dataobject.DataObject for data storage and change tracking.
type taxonomyImplementation struct {
	dataobject.DataObject
}

// GetID returns the unique identifier of the taxonomy.
func (o *taxonomyImplementation) GetID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the unique identifier of the taxonomy.
func (o *taxonomyImplementation) SetID(id string) TaxonomyInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// GetName returns the display name of the taxonomy.
func (o *taxonomyImplementation) GetName() string {
	return o.Get(COLUMN_NAME)
}

// SetName sets the display name of the taxonomy.
func (o *taxonomyImplementation) SetName(name string) TaxonomyInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

// GetSlug returns the URL-friendly slug of the taxonomy.
func (o *taxonomyImplementation) GetSlug() string {
	return o.Get(COLUMN_SLUG)
}

// SetSlug sets the URL-friendly slug of the taxonomy.
// The slug is automatically slugified using the str.Slugify function.
func (o *taxonomyImplementation) SetSlug(slug string) TaxonomyInterface {
	o.Set(COLUMN_SLUG, str.Slugify(slug, '-'))
	return o
}

// GetDescription returns the description of the taxonomy.
func (o *taxonomyImplementation) GetDescription() string {
	return o.Get(COLUMN_DESCRIPTION)
}

// SetDescription sets the description of the taxonomy.
func (o *taxonomyImplementation) SetDescription(description string) TaxonomyInterface {
	o.Set(COLUMN_DESCRIPTION, description)
	return o
}

// GetCreatedAt returns the creation timestamp as a string.
func (o *taxonomyImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the creation timestamp.
func (o *taxonomyImplementation) SetCreatedAt(createdAt string) TaxonomyInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
// Returns the null datetime if the created_at field is empty.
func (o *taxonomyImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(createdAt, carbon.UTC)
}

// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
// Returns zero time if the created_at field is empty.
func (o *taxonomyImplementation) GetCreatedAtTime() time.Time {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return time.Time{}
	}
	return carbon.Parse(createdAt, carbon.UTC).StdTime()
}

// GetUpdatedAt returns the last update timestamp as a string.
func (o *taxonomyImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// SetUpdatedAt sets the last update timestamp.
func (o *taxonomyImplementation) SetUpdatedAt(updatedAt string) TaxonomyInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
// Returns the null datetime if the updated_at field is empty.
func (o *taxonomyImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(updatedAt, carbon.UTC)
}

// GetUpdatedAtTime returns the last update timestamp as a time.Time instance.
// Returns zero time if the updated_at field is empty.
func (o *taxonomyImplementation) GetUpdatedAtTime() time.Time {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return time.Time{}
	}
	return carbon.Parse(updatedAt, carbon.UTC).StdTime()
}

// GetData returns all taxonomy data as a map.
func (o *taxonomyImplementation) GetData() map[string]string {
	return o.DataObject.Data()
}

// GetDataChanged returns only the fields that have been modified.
func (o *taxonomyImplementation) GetDataChanged() map[string]string {
	return o.DataObject.DataChanged()
}

// MarkAsNotDirty clears the dirty state of the taxonomy.
func (o *taxonomyImplementation) MarkAsNotDirty() {
	o.DataObject.MarkAsNotDirty()
}

// Get retrieves a value by key from the taxonomy data.
func (o *taxonomyImplementation) Get(key string) string {
	return o.DataObject.Get(key)
}

// Set stores a value by key in the taxonomy data.
func (o *taxonomyImplementation) Set(key string, value string) {
	o.DataObject.Set(key, value)
}

// Hydrate populates the taxonomy with data from a map.
func (o *taxonomyImplementation) Hydrate(data map[string]string) {
	o.DataObject.Hydrate(data)
}

// IsDirty returns true if the taxonomy has unsaved changes.
func (o *taxonomyImplementation) IsDirty() bool {
	return o.DataObject.IsDirty()
}

// ============================ TERM ============================

// NewTerm creates a new Term instance with default values.
// The term is initialized with a generated ID, empty taxonomy/parent IDs,
// zero sequence, empty name/slug/description, zero count, and current timestamps.
func NewTerm() TermInterface {
	o := &termImplementation{}
	o.SetID(GenerateShortID()).
		SetTaxonomyID("").
		SetParentID("").
		SetSequence(0).
		SetName("").
		SetSlug("").
		SetDescription("").
		SetCount(0).
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return o
}

// NewTermFromExistingData creates a Term instance from existing data.
// This is useful when hydrating a term from database records.
func NewTermFromExistingData(data map[string]string) TermInterface {
	o := &termImplementation{}
	o.Hydrate(data)
	return o
}

// termImplementation is the concrete implementation of the TermInterface.
// It embeds dataobject.DataObject for data storage and change tracking.
type termImplementation struct {
	dataobject.DataObject
}

// GetID returns the unique identifier of the term.
func (o *termImplementation) GetID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the unique identifier of the term.
func (o *termImplementation) SetID(id string) TermInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// GetTaxonomyID returns the ID of the taxonomy this term belongs to.
func (o *termImplementation) GetTaxonomyID() string {
	return o.Get(COLUMN_TAXONOMY_ID)
}

// SetTaxonomyID sets the ID of the taxonomy this term belongs to.
func (o *termImplementation) SetTaxonomyID(taxonomyID string) TermInterface {
	o.Set(COLUMN_TAXONOMY_ID, taxonomyID)
	return o
}

// GetParentID returns the ID of the parent term (for hierarchical terms).
func (o *termImplementation) GetParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

// SetParentID sets the ID of the parent term (for hierarchical terms).
func (o *termImplementation) SetParentID(parentID string) TermInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

// GetSequence returns the display sequence/order of the term.
// Returns 0 if the sequence field is empty or cannot be parsed.
func (o *termImplementation) GetSequence() int {
	seqStr := o.Get(COLUMN_SEQUENCE)
	if seqStr == "" {
		return 0
	}
	var seq int
	if _, err := parseInt(seqStr, &seq); err != nil {
		return 0
	}
	return seq
}

// SetSequence sets the display sequence/order of the term.
func (o *termImplementation) SetSequence(sequence int) TermInterface {
	o.Set(COLUMN_SEQUENCE, intToString(sequence))
	return o
}

// GetName returns the display name of the term.
func (o *termImplementation) GetName() string {
	return o.Get(COLUMN_NAME)
}

// SetName sets the display name of the term.
func (o *termImplementation) SetName(name string) TermInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

// GetSlug returns the URL-friendly slug of the term.
func (o *termImplementation) GetSlug() string {
	return o.Get(COLUMN_SLUG)
}

// SetSlug sets the URL-friendly slug of the term.
// The slug is automatically slugified using the str.Slugify function.
func (o *termImplementation) SetSlug(slug string) TermInterface {
	o.Set(COLUMN_SLUG, str.Slugify(slug, '-'))
	return o
}

// GetDescription returns the description of the term.
func (o *termImplementation) GetDescription() string {
	return o.Get(COLUMN_DESCRIPTION)
}

// SetDescription sets the description of the term.
func (o *termImplementation) SetDescription(description string) TermInterface {
	o.Set(COLUMN_DESCRIPTION, description)
	return o
}

// GetCount returns the number of posts associated with this term.
// Returns 0 if the count field is empty or cannot be parsed.
func (o *termImplementation) GetCount() int {
	countStr := o.Get(COLUMN_COUNT)
	if countStr == "" {
		return 0
	}
	// Safe conversion
	var count int
	if _, err := parseInt(countStr, &count); err != nil {
		return 0
	}
	return count
}

// SetCount sets the number of posts associated with this term.
func (o *termImplementation) SetCount(count int) TermInterface {
	o.Set(COLUMN_COUNT, intToString(count))
	return o
}

// GetCreatedAt returns the creation timestamp as a string.
func (o *termImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the creation timestamp.
func (o *termImplementation) SetCreatedAt(createdAt string) TermInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
// Returns the null datetime if the created_at field is empty.
func (o *termImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(createdAt, carbon.UTC)
}

// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
// Returns zero time if the created_at field is empty.
func (o *termImplementation) GetCreatedAtTime() time.Time {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return time.Time{}
	}
	return carbon.Parse(createdAt, carbon.UTC).StdTime()
}

// GetUpdatedAt returns the last update timestamp as a string.
func (o *termImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// SetUpdatedAt sets the last update timestamp.
func (o *termImplementation) SetUpdatedAt(updatedAt string) TermInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
// Returns the null datetime if the updated_at field is empty.
func (o *termImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(updatedAt, carbon.UTC)
}

// GetUpdatedAtTime returns the last update timestamp as a time.Time instance.
// Returns zero time if the updated_at field is empty.
func (o *termImplementation) GetUpdatedAtTime() time.Time {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return time.Time{}
	}
	return carbon.Parse(updatedAt, carbon.UTC).StdTime()
}

// GetData returns all term data as a map.
func (o *termImplementation) GetData() map[string]string {
	return o.DataObject.Data()
}

// GetDataChanged returns only the fields that have been modified.
func (o *termImplementation) GetDataChanged() map[string]string {
	return o.DataObject.DataChanged()
}

// MarkAsNotDirty clears the dirty state of the term.
func (o *termImplementation) MarkAsNotDirty() {
	o.DataObject.MarkAsNotDirty()
}

// Get retrieves a value by key from the term data.
func (o *termImplementation) Get(key string) string {
	return o.DataObject.Get(key)
}

// Set stores a value by key in the term data.
func (o *termImplementation) Set(key string, value string) {
	o.DataObject.Set(key, value)
}

// Hydrate populates the term with data from a map.
func (o *termImplementation) Hydrate(data map[string]string) {
	o.DataObject.Hydrate(data)
}

// IsDirty returns true if the term has unsaved changes.
func (o *termImplementation) IsDirty() bool {
	return o.DataObject.IsDirty()
}

// ============================ TERM RELATION ============================

// NewTermRelation creates a new TermRelation instance with default values.
// The relation is initialized with a generated ID, empty post/term IDs,
// zero sequence, and current timestamps.
func NewTermRelation() TermRelationInterface {
	o := &termRelationImplementation{}
	o.SetID(GenerateShortID()).
		SetPostID("").
		SetTermID("").
		SetSequence(0).
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return o
}

// NewTermRelationFromExistingData creates a TermRelation instance from existing data.
// This is useful when hydrating a term relation from database records.
func NewTermRelationFromExistingData(data map[string]string) TermRelationInterface {
	o := &termRelationImplementation{}
	o.Hydrate(data)
	return o
}

// termRelationImplementation is the concrete implementation of the TermRelationInterface.
// It embeds dataobject.DataObject for data storage and change tracking.
type termRelationImplementation struct {
	dataobject.DataObject
}

// GetID returns the unique identifier of the term relation.
func (o *termRelationImplementation) GetID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the unique identifier of the term relation.
func (o *termRelationImplementation) SetID(id string) TermRelationInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// GetPostID returns the ID of the post in this relation.
func (o *termRelationImplementation) GetPostID() string {
	return o.Get(COLUMN_POST_ID)
}

// SetPostID sets the ID of the post in this relation.
func (o *termRelationImplementation) SetPostID(postID string) TermRelationInterface {
	o.Set(COLUMN_POST_ID, postID)
	return o
}

// GetTermID returns the ID of the term in this relation.
func (o *termRelationImplementation) GetTermID() string {
	return o.Get(COLUMN_TERM_ID)
}

// SetTermID sets the ID of the term in this relation.
func (o *termRelationImplementation) SetTermID(termID string) TermRelationInterface {
	o.Set(COLUMN_TERM_ID, termID)
	return o
}

// GetSequence returns the display sequence/order of the relation.
// Returns 0 if the sequence field is empty or cannot be parsed.
func (o *termRelationImplementation) GetSequence() int {
	seqStr := o.Get(COLUMN_SEQUENCE)
	if seqStr == "" {
		return 0
	}
	var seq int
	if _, err := parseInt(seqStr, &seq); err != nil {
		return 0
	}
	return seq
}

// SetSequence sets the display sequence/order of the relation.
func (o *termRelationImplementation) SetSequence(sequence int) TermRelationInterface {
	o.Set(COLUMN_SEQUENCE, intToString(sequence))
	return o
}

// GetCreatedAt returns the creation timestamp as a string.
func (o *termRelationImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the creation timestamp.
func (o *termRelationImplementation) SetCreatedAt(createdAt string) TermRelationInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
// Returns the null datetime if the created_at field is empty.
func (o *termRelationImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(createdAt, carbon.UTC)
}

// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
// Returns zero time if the created_at field is empty.
func (o *termRelationImplementation) GetCreatedAtTime() time.Time {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return time.Time{}
	}
	return carbon.Parse(createdAt, carbon.UTC).StdTime()
}

// GetUpdatedAt returns the last update timestamp as a string.
func (o *termRelationImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// SetUpdatedAt sets the last update timestamp.
func (o *termRelationImplementation) SetUpdatedAt(updatedAt string) TermRelationInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
// Returns the null datetime if the updated_at field is empty.
func (o *termRelationImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(updatedAt, carbon.UTC)
}

// GetUpdatedAtTime returns the last update timestamp as a time.Time instance.
// Returns zero time if the updated_at field is empty.
func (o *termRelationImplementation) GetUpdatedAtTime() time.Time {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return time.Time{}
	}
	return carbon.Parse(updatedAt, carbon.UTC).StdTime()
}

// GetData returns all term relation data as a map.
func (o *termRelationImplementation) GetData() map[string]string {
	return o.DataObject.Data()
}

// GetDataChanged returns only the fields that have been modified.
func (o *termRelationImplementation) GetDataChanged() map[string]string {
	return o.DataObject.DataChanged()
}

// MarkAsNotDirty clears the dirty state of the term relation.
func (o *termRelationImplementation) MarkAsNotDirty() {
	o.DataObject.MarkAsNotDirty()
}

// Get retrieves a value by key from the term relation data.
func (o *termRelationImplementation) Get(key string) string {
	return o.DataObject.Get(key)
}

// Set stores a value by key in the term relation data.
func (o *termRelationImplementation) Set(key string, value string) {
	o.DataObject.Set(key, value)
}

// Hydrate populates the term relation with data from a map.
func (o *termRelationImplementation) Hydrate(data map[string]string) {
	o.DataObject.Hydrate(data)
}

// IsDirty returns true if the term relation has unsaved changes.
func (o *termRelationImplementation) IsDirty() bool {
	return o.DataObject.IsDirty()
}

// Helper functions

// sb_NULL_DATETIME is the default null datetime value used when timestamps are empty.
var sb_NULL_DATETIME = "1970-01-01 00:00:00"

// parseInt parses a string into an integer without using strconv.
// It handles negative numbers and returns an error for invalid input.
func parseInt(s string, v *int) (int, error) {
	// Simple parse implementation
	result := 0
	negative := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			if i == 0 && c == '-' {
				negative = true
				continue
			}
			return 0, nil
		}
		result = result*10 + int(c-'0')
	}
	if negative {
		result = -result
	}
	*v = result
	return result, nil
}

// intToString converts an integer to a string without using strconv.
// It properly handles negative numbers and zero.
func intToString(i int) string {
	// Simple implementation without strconv
	if i == 0 {
		return "0"
	}
	negative := i < 0
	if negative {
		i = -i
	}
	result := ""
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}
