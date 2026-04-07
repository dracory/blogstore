package blogstore

import (
	"time"

	"github.com/dracory/dataobject"
	"github.com/dracory/str"
	"github.com/dromara/carbon/v2"
)

// ============================ TAXONOMY ============================

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

func NewTaxonomyFromExistingData(data map[string]string) TaxonomyInterface {
	o := &taxonomyImplementation{}
	o.Hydrate(data)
	return o
}

type taxonomyImplementation struct {
	dataobject.DataObject
}

func (o *taxonomyImplementation) GetID() string {
	return o.Get(COLUMN_ID)
}

func (o *taxonomyImplementation) SetID(id string) TaxonomyInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *taxonomyImplementation) GetName() string {
	return o.Get(COLUMN_NAME)
}

func (o *taxonomyImplementation) SetName(name string) TaxonomyInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *taxonomyImplementation) GetSlug() string {
	return o.Get(COLUMN_SLUG)
}

func (o *taxonomyImplementation) SetSlug(slug string) TaxonomyInterface {
	o.Set(COLUMN_SLUG, str.Slugify(slug, '-'))
	return o
}

func (o *taxonomyImplementation) GetDescription() string {
	return o.Get(COLUMN_DESCRIPTION)
}

func (o *taxonomyImplementation) SetDescription(description string) TaxonomyInterface {
	o.Set(COLUMN_DESCRIPTION, description)
	return o
}

func (o *taxonomyImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *taxonomyImplementation) SetCreatedAt(createdAt string) TaxonomyInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *taxonomyImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(createdAt, carbon.UTC)
}

func (o *taxonomyImplementation) GetCreatedAtTime() time.Time {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return time.Time{}
	}
	return carbon.Parse(createdAt, carbon.UTC).StdTime()
}

func (o *taxonomyImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *taxonomyImplementation) SetUpdatedAt(updatedAt string) TaxonomyInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *taxonomyImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(updatedAt, carbon.UTC)
}

func (o *taxonomyImplementation) GetUpdatedAtTime() time.Time {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return time.Time{}
	}
	return carbon.Parse(updatedAt, carbon.UTC).StdTime()
}

func (o *taxonomyImplementation) GetData() map[string]string {
	return o.DataObject.Data()
}

func (o *taxonomyImplementation) GetDataChanged() map[string]string {
	return o.DataObject.DataChanged()
}

func (o *taxonomyImplementation) MarkAsNotDirty() {
	o.DataObject.MarkAsNotDirty()
}

func (o *taxonomyImplementation) Get(key string) string {
	return o.DataObject.Get(key)
}

func (o *taxonomyImplementation) Set(key string, value string) {
	o.DataObject.Set(key, value)
}

func (o *taxonomyImplementation) Hydrate(data map[string]string) {
	o.DataObject.Hydrate(data)
}

func (o *taxonomyImplementation) IsDirty() bool {
	return o.DataObject.IsDirty()
}

// ============================ TERM ============================

func NewTerm() TermInterface {
	o := &termImplementation{}
	o.SetID(GenerateShortID()).
		SetTaxonomyID("").
		SetParentID("").
		SetName("").
		SetSlug("").
		SetDescription("").
		SetCount(0).
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return o
}

func NewTermFromExistingData(data map[string]string) TermInterface {
	o := &termImplementation{}
	o.Hydrate(data)
	return o
}

type termImplementation struct {
	dataobject.DataObject
}

func (o *termImplementation) GetID() string {
	return o.Get(COLUMN_ID)
}

func (o *termImplementation) SetID(id string) TermInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *termImplementation) GetTaxonomyID() string {
	return o.Get(COLUMN_TAXONOMY_ID)
}

func (o *termImplementation) SetTaxonomyID(taxonomyID string) TermInterface {
	o.Set(COLUMN_TAXONOMY_ID, taxonomyID)
	return o
}

func (o *termImplementation) GetParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

func (o *termImplementation) SetParentID(parentID string) TermInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

func (o *termImplementation) GetName() string {
	return o.Get(COLUMN_NAME)
}

func (o *termImplementation) SetName(name string) TermInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *termImplementation) GetSlug() string {
	return o.Get(COLUMN_SLUG)
}

func (o *termImplementation) SetSlug(slug string) TermInterface {
	o.Set(COLUMN_SLUG, str.Slugify(slug, '-'))
	return o
}

func (o *termImplementation) GetDescription() string {
	return o.Get(COLUMN_DESCRIPTION)
}

func (o *termImplementation) SetDescription(description string) TermInterface {
	o.Set(COLUMN_DESCRIPTION, description)
	return o
}

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

func (o *termImplementation) SetCount(count int) TermInterface {
	o.Set(COLUMN_COUNT, intToString(count))
	return o
}

func (o *termImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *termImplementation) SetCreatedAt(createdAt string) TermInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *termImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(createdAt, carbon.UTC)
}

func (o *termImplementation) GetCreatedAtTime() time.Time {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return time.Time{}
	}
	return carbon.Parse(createdAt, carbon.UTC).StdTime()
}

func (o *termImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *termImplementation) SetUpdatedAt(updatedAt string) TermInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *termImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(updatedAt, carbon.UTC)
}

func (o *termImplementation) GetUpdatedAtTime() time.Time {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return time.Time{}
	}
	return carbon.Parse(updatedAt, carbon.UTC).StdTime()
}

func (o *termImplementation) GetData() map[string]string {
	return o.DataObject.Data()
}

func (o *termImplementation) GetDataChanged() map[string]string {
	return o.DataObject.DataChanged()
}

func (o *termImplementation) MarkAsNotDirty() {
	o.DataObject.MarkAsNotDirty()
}

func (o *termImplementation) Get(key string) string {
	return o.DataObject.Get(key)
}

func (o *termImplementation) Set(key string, value string) {
	o.DataObject.Set(key, value)
}

func (o *termImplementation) Hydrate(data map[string]string) {
	o.DataObject.Hydrate(data)
}

func (o *termImplementation) IsDirty() bool {
	return o.DataObject.IsDirty()
}

// ============================ TERM RELATION ============================

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

func NewTermRelationFromExistingData(data map[string]string) TermRelationInterface {
	o := &termRelationImplementation{}
	o.Hydrate(data)
	return o
}

type termRelationImplementation struct {
	dataobject.DataObject
}

func (o *termRelationImplementation) GetID() string {
	return o.Get(COLUMN_ID)
}

func (o *termRelationImplementation) SetID(id string) TermRelationInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *termRelationImplementation) GetPostID() string {
	return o.Get(COLUMN_POST_ID)
}

func (o *termRelationImplementation) SetPostID(postID string) TermRelationInterface {
	o.Set(COLUMN_POST_ID, postID)
	return o
}

func (o *termRelationImplementation) GetTermID() string {
	return o.Get(COLUMN_TERM_ID)
}

func (o *termRelationImplementation) SetTermID(termID string) TermRelationInterface {
	o.Set(COLUMN_TERM_ID, termID)
	return o
}

func (o *termRelationImplementation) GetSequence() int {
	seqStr := o.Get(COLUMN_TERM_SEQUENCE)
	if seqStr == "" {
		return 0
	}
	var seq int
	if _, err := parseInt(seqStr, &seq); err != nil {
		return 0
	}
	return seq
}

func (o *termRelationImplementation) SetSequence(sequence int) TermRelationInterface {
	o.Set(COLUMN_TERM_SEQUENCE, intToString(sequence))
	return o
}

func (o *termRelationImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *termRelationImplementation) SetCreatedAt(createdAt string) TermRelationInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *termRelationImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(createdAt, carbon.UTC)
}

func (o *termRelationImplementation) GetCreatedAtTime() time.Time {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return time.Time{}
	}
	return carbon.Parse(createdAt, carbon.UTC).StdTime()
}

func (o *termRelationImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *termRelationImplementation) SetUpdatedAt(updatedAt string) TermRelationInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *termRelationImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return carbon.Parse(sb_NULL_DATETIME, carbon.UTC)
	}
	return carbon.Parse(updatedAt, carbon.UTC)
}

func (o *termRelationImplementation) GetUpdatedAtTime() time.Time {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return time.Time{}
	}
	return carbon.Parse(updatedAt, carbon.UTC).StdTime()
}

func (o *termRelationImplementation) GetData() map[string]string {
	return o.DataObject.Data()
}

func (o *termRelationImplementation) GetDataChanged() map[string]string {
	return o.DataObject.DataChanged()
}

func (o *termRelationImplementation) MarkAsNotDirty() {
	o.DataObject.MarkAsNotDirty()
}

func (o *termRelationImplementation) Get(key string) string {
	return o.DataObject.Get(key)
}

func (o *termRelationImplementation) Set(key string, value string) {
	o.DataObject.Set(key, value)
}

func (o *termRelationImplementation) Hydrate(data map[string]string) {
	o.DataObject.Hydrate(data)
}

func (o *termRelationImplementation) IsDirty() bool {
	return o.DataObject.IsDirty()
}

// Helper functions
var sb_NULL_DATETIME = "1970-01-01 00:00:00"

func parseInt(s string, v *int) (int, error) {
	// Simple parse implementation
	result := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			if i == 0 && c == '-' {
				continue
			}
			return 0, nil
		}
		result = result*10 + int(c-'0')
	}
	*v = result
	return result, nil
}

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
