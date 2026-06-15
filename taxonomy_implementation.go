package blogstore

import (
	"strconv"
	"time"

	"github.com/dracory/neat/database/orm"
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
// It uses neat ORM traits for data storage.
type taxonomyImplementation struct {
	orm.ShortID

	NameField        string `db:"name"`
	SlugField        string `db:"slug"`
	DescriptionField string `db:"description"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
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
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

// SetCreatedAt sets the creation timestamp.
func (o *taxonomyImplementation) SetCreatedAt(createdAt string) TaxonomyInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
// Returns the null datetime if the created_at field is empty.
func (o *taxonomyImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
// Returns zero time if the created_at field is empty.
func (o *taxonomyImplementation) GetCreatedAtTime() time.Time {
	return o.CreatedAtField.CreatedAt
}

// GetUpdatedAt returns the last update timestamp as a string.
func (o *taxonomyImplementation) GetUpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

// SetUpdatedAt sets the last update timestamp.
func (o *taxonomyImplementation) SetUpdatedAt(updatedAt string) TaxonomyInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
// Returns the null datetime if the updated_at field is empty.
func (o *taxonomyImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

// GetUpdatedAtTime returns the last update timestamp as a time.Time instance.
// Returns zero time if the updated_at field is empty.
func (o *taxonomyImplementation) GetUpdatedAtTime() time.Time {
	return o.UpdatedAtField.UpdatedAt
}

// GetData returns all taxonomy data as a map.
func (o *taxonomyImplementation) GetData() map[string]string {
	var createdAt, updatedAt string
	if !o.CreatedAtField.CreatedAt.IsZero() {
		createdAt = carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString(carbon.UTC)
	}
	if !o.UpdatedAtField.UpdatedAt.IsZero() {
		updatedAt = carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString(carbon.UTC)
	}

	return map[string]string{
		COLUMN_ID:          o.ShortID.ID,
		COLUMN_NAME:        o.NameField,
		COLUMN_SLUG:        o.SlugField,
		COLUMN_DESCRIPTION: o.DescriptionField,
		COLUMN_CREATED_AT:  createdAt,
		COLUMN_UPDATED_AT:  updatedAt,
	}
}

// GetDataChanged returns only the fields that have been modified.
// Since neat ORM traits don't track dirty state, return all fields as changed.
func (o *taxonomyImplementation) GetDataChanged() map[string]string {
	return o.GetData()
}

// MarkAsNotDirty clears the dirty state of the taxonomy.
// No-op since neat ORM traits don't track dirty state.
func (o *taxonomyImplementation) MarkAsNotDirty(columns ...string) {
}

// MarkAsDirty marks the taxonomy as dirty.
// No-op since neat ORM traits don't track dirty state.
func (o *taxonomyImplementation) MarkAsDirty(columns ...string) {
}

// Get retrieves a value by key from the taxonomy data.
func (o *taxonomyImplementation) Get(key string) string {
	switch key {
	case COLUMN_ID:
		return o.ID
	case COLUMN_NAME:
		return o.NameField
	case COLUMN_SLUG:
		return o.SlugField
	case COLUMN_DESCRIPTION:
		return o.DescriptionField
	case COLUMN_CREATED_AT:
		if o.CreatedAtField.CreatedAt.IsZero() {
			return ""
		}
		return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString(carbon.UTC)
	case COLUMN_UPDATED_AT:
		if o.UpdatedAtField.UpdatedAt.IsZero() {
			return ""
		}
		return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString(carbon.UTC)
	default:
		return ""
	}
}

// Set stores a value by key in the taxonomy data.
func (o *taxonomyImplementation) Set(key string, value string) {
	switch key {
	case COLUMN_ID:
		o.ShortID.ID = value
	case COLUMN_NAME:
		o.NameField = value
	case COLUMN_SLUG:
		o.SlugField = value
	case COLUMN_DESCRIPTION:
		o.DescriptionField = value
	case COLUMN_CREATED_AT:
		if value != "" {
			o.CreatedAtField.CreatedAt = carbon.Parse(value, carbon.UTC).StdTime()
		}
	case COLUMN_UPDATED_AT:
		if value != "" {
			o.UpdatedAtField.UpdatedAt = carbon.Parse(value, carbon.UTC).StdTime()
		}
	}
}

// Hydrate populates the taxonomy with data from a map.
func (o *taxonomyImplementation) Hydrate(data map[string]string) {
	for key, value := range data {
		o.Set(key, value)
	}
}

// IsDirty returns true if the taxonomy has unsaved changes.
// Always returns false since neat ORM traits don't track dirty state.
func (o *taxonomyImplementation) IsDirty() bool {
	return false
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
// It uses neat ORM traits for data storage.
type termImplementation struct {
	orm.ShortID

	TaxonomyIDField  string `db:"taxonomy_id"`
	ParentIDField    string `db:"parent_id"`
	SequenceField    int    `db:"sequence"`
	NameField        string `db:"name"`
	SlugField        string `db:"slug"`
	DescriptionField string `db:"description"`
	CountField       int    `db:"count"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
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
func (o *termImplementation) GetSequence() int {
	return o.SequenceField
}

// SetSequence sets the display sequence/order of the term.
func (o *termImplementation) SetSequence(sequence int) TermInterface {
	o.SequenceField = sequence
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
func (o *termImplementation) GetCount() int {
	return o.CountField
}

// SetCount sets the number of posts associated with this term.
func (o *termImplementation) SetCount(count int) TermInterface {
	o.CountField = count
	return o
}

// GetCreatedAt returns the creation timestamp as a string.
func (o *termImplementation) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

// SetCreatedAt sets the creation timestamp.
func (o *termImplementation) SetCreatedAt(createdAt string) TermInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
// Returns the null datetime if the created_at field is empty.
func (o *termImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
// Returns zero time if the created_at field is empty.
func (o *termImplementation) GetCreatedAtTime() time.Time {
	return o.CreatedAtField.CreatedAt
}

// GetUpdatedAt returns the last update timestamp as a string.
func (o *termImplementation) GetUpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

// SetUpdatedAt sets the last update timestamp.
func (o *termImplementation) SetUpdatedAt(updatedAt string) TermInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
// Returns the null datetime if the updated_at field is empty.
func (o *termImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

// GetUpdatedAtTime returns the last update timestamp as a time.Time instance.
// Returns zero time if the updated_at field is empty.
func (o *termImplementation) GetUpdatedAtTime() time.Time {
	return o.UpdatedAtField.UpdatedAt
}

// GetData returns all term data as a map.
func (o *termImplementation) GetData() map[string]string {
	var createdAt, updatedAt string
	if !o.CreatedAtField.CreatedAt.IsZero() {
		createdAt = carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString(carbon.UTC)
	}
	if !o.UpdatedAtField.UpdatedAt.IsZero() {
		updatedAt = carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString(carbon.UTC)
	}

	return map[string]string{
		COLUMN_ID:          o.ShortID.ID,
		COLUMN_TAXONOMY_ID: o.TaxonomyIDField,
		COLUMN_PARENT_ID:   o.ParentIDField,
		COLUMN_SEQUENCE:    strconv.Itoa(o.SequenceField),
		COLUMN_NAME:        o.NameField,
		COLUMN_SLUG:        o.SlugField,
		COLUMN_DESCRIPTION: o.DescriptionField,
		COLUMN_COUNT:       strconv.Itoa(o.CountField),
		COLUMN_CREATED_AT:  createdAt,
		COLUMN_UPDATED_AT:  updatedAt,
	}
}

// GetDataChanged returns only the fields that have been modified.
// Since neat ORM traits don't track dirty state, return all fields as changed.
func (o *termImplementation) GetDataChanged() map[string]string {
	return o.GetData()
}

// MarkAsNotDirty clears the dirty state of the term.
// No-op since neat ORM traits don't track dirty state.
func (o *termImplementation) MarkAsNotDirty(columns ...string) {
}

// MarkAsDirty marks the term as dirty.
// No-op since neat ORM traits don't track dirty state.
func (o *termImplementation) MarkAsDirty(columns ...string) {
}

// Get retrieves a value by key from the term data.
func (o *termImplementation) Get(key string) string {
	switch key {
	case COLUMN_ID:
		return o.ID
	case COLUMN_TAXONOMY_ID:
		return o.TaxonomyIDField
	case COLUMN_PARENT_ID:
		return o.ParentIDField
	case COLUMN_SEQUENCE:
		return strconv.Itoa(o.SequenceField)
	case COLUMN_NAME:
		return o.NameField
	case COLUMN_SLUG:
		return o.SlugField
	case COLUMN_DESCRIPTION:
		return o.DescriptionField
	case COLUMN_COUNT:
		return strconv.Itoa(o.CountField)
	case COLUMN_CREATED_AT:
		if o.CreatedAtField.CreatedAt.IsZero() {
			return ""
		}
		return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString(carbon.UTC)
	case COLUMN_UPDATED_AT:
		if o.UpdatedAtField.UpdatedAt.IsZero() {
			return ""
		}
		return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString(carbon.UTC)
	default:
		return ""
	}
}

// Set stores a value by key in the term data.
func (o *termImplementation) Set(key string, value string) {
	switch key {
	case COLUMN_ID:
		o.ShortID.ID = value
	case COLUMN_TAXONOMY_ID:
		o.TaxonomyIDField = value
	case COLUMN_PARENT_ID:
		o.ParentIDField = value
	case COLUMN_SEQUENCE:
		if seq, err := strconv.Atoi(value); err == nil {
			o.SequenceField = seq
		}
	case COLUMN_NAME:
		o.NameField = value
	case COLUMN_SLUG:
		o.SlugField = value
	case COLUMN_DESCRIPTION:
		o.DescriptionField = value
	case COLUMN_COUNT:
		if count, err := strconv.Atoi(value); err == nil {
			o.CountField = count
		}
	case COLUMN_CREATED_AT:
		if value != "" {
			o.CreatedAtField.CreatedAt = carbon.Parse(value, carbon.UTC).StdTime()
		}
	case COLUMN_UPDATED_AT:
		if value != "" {
			o.UpdatedAtField.UpdatedAt = carbon.Parse(value, carbon.UTC).StdTime()
		}
	}
}

// Hydrate populates the term with data from a map.
func (o *termImplementation) Hydrate(data map[string]string) {
	for key, value := range data {
		o.Set(key, value)
	}
}

// IsDirty returns true if the term has unsaved changes.
// Always returns false since neat ORM traits don't track dirty state.
func (o *termImplementation) IsDirty() bool {
	return false
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
// It uses neat ORM traits for data storage.
type termRelationImplementation struct {
	orm.ShortID

	PostIDField   string `db:"post_id"`
	TermIDField   string `db:"term_id"`
	SequenceField int    `db:"sequence"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
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
func (o *termRelationImplementation) GetSequence() int {
	return o.SequenceField
}

// SetSequence sets the display sequence/order of the relation.
func (o *termRelationImplementation) SetSequence(sequence int) TermRelationInterface {
	o.SequenceField = sequence
	return o
}

// GetCreatedAt returns the creation timestamp as a string.
func (o *termRelationImplementation) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

// SetCreatedAt sets the creation timestamp.
func (o *termRelationImplementation) SetCreatedAt(createdAt string) TermRelationInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
// Returns the null datetime if the created_at field is empty.
func (o *termRelationImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
// Returns zero time if the created_at field is empty.
func (o *termRelationImplementation) GetCreatedAtTime() time.Time {
	return o.CreatedAtField.CreatedAt
}

// GetUpdatedAt returns the last update timestamp as a string.
func (o *termRelationImplementation) GetUpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

// SetUpdatedAt sets the last update timestamp.
func (o *termRelationImplementation) SetUpdatedAt(updatedAt string) TermRelationInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
// Returns the null datetime if the updated_at field is empty.
func (o *termRelationImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

// GetUpdatedAtTime returns the last update timestamp as a time.Time instance.
// Returns zero time if the updated_at field is empty.
func (o *termRelationImplementation) GetUpdatedAtTime() time.Time {
	return o.UpdatedAtField.UpdatedAt
}

// GetData returns all term relation data as a map.
func (o *termRelationImplementation) GetData() map[string]string {
	var createdAt, updatedAt string
	if !o.CreatedAtField.CreatedAt.IsZero() {
		createdAt = carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString(carbon.UTC)
	}
	if !o.UpdatedAtField.UpdatedAt.IsZero() {
		updatedAt = carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString(carbon.UTC)
	}

	return map[string]string{
		COLUMN_ID:         o.ShortID.ID,
		COLUMN_POST_ID:    o.PostIDField,
		COLUMN_TERM_ID:    o.TermIDField,
		COLUMN_SEQUENCE:   strconv.Itoa(o.SequenceField),
		COLUMN_CREATED_AT: createdAt,
		COLUMN_UPDATED_AT: updatedAt,
	}
}

// GetDataChanged returns only the fields that have been modified.
// Since neat ORM traits don't track dirty state, return all fields as changed.
func (o *termRelationImplementation) GetDataChanged() map[string]string {
	return o.GetData()
}

// MarkAsNotDirty clears the dirty state of the term relation.
// No-op since neat ORM traits don't track dirty state.
func (o *termRelationImplementation) MarkAsNotDirty(columns ...string) {
}

// MarkAsDirty marks the term relation as dirty.
// No-op since neat ORM traits don't track dirty state.
func (o *termRelationImplementation) MarkAsDirty(columns ...string) {
}

// Get retrieves a value by key from the term relation data.
func (o *termRelationImplementation) Get(key string) string {
	switch key {
	case COLUMN_ID:
		return o.ShortID.ID
	case COLUMN_POST_ID:
		return o.PostIDField
	case COLUMN_TERM_ID:
		return o.TermIDField
	case COLUMN_SEQUENCE:
		return strconv.Itoa(o.SequenceField)
	case COLUMN_CREATED_AT:
		if o.CreatedAtField.CreatedAt.IsZero() {
			return ""
		}
		return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString(carbon.UTC)
	case COLUMN_UPDATED_AT:
		if o.UpdatedAtField.UpdatedAt.IsZero() {
			return ""
		}
		return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString(carbon.UTC)
	default:
		return ""
	}
}

// Set stores a value by key in the term relation data.
func (o *termRelationImplementation) Set(key string, value string) {
	switch key {
	case COLUMN_ID:
		o.ShortID.ID = value
	case COLUMN_POST_ID:
		o.PostIDField = value
	case COLUMN_TERM_ID:
		o.TermIDField = value
	case COLUMN_SEQUENCE:
		if seq, err := strconv.Atoi(value); err == nil {
			o.SequenceField = seq
		}
	case COLUMN_CREATED_AT:
		if value != "" {
			o.CreatedAtField.CreatedAt = carbon.Parse(value, carbon.UTC).StdTime()
		}
	case COLUMN_UPDATED_AT:
		if value != "" {
			o.UpdatedAtField.UpdatedAt = carbon.Parse(value, carbon.UTC).StdTime()
		}
	}
}

// Hydrate populates the term relation with data from a map.
func (o *termRelationImplementation) Hydrate(data map[string]string) {
	for key, value := range data {
		o.Set(key, value)
	}
}

// IsDirty returns true if the term relation has unsaved changes.
// Always returns false since neat ORM traits don't track dirty state.
func (o *termRelationImplementation) IsDirty() bool {
	return false
}
