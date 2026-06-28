package blogstore

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	"github.com/dromara/carbon/v2"
)

// MediaInterface defines the interface for media (file attachment) operations.
// Media files represent files attached to any entity (post, category, etc),
// such as images, videos, or documents.
type MediaInterface interface {
	// IsSoftDeleted returns true if the media is soft deleted.
	IsSoftDeleted() bool

	// GetID returns the unique identifier of the media.
	GetID() string
	// SetID sets the unique identifier of the media.
	SetID(id string) MediaInterface

	// GetEntityID returns the ID of the entity this media is attached to.
	GetEntityID() string
	// SetEntityID sets the ID of the entity this media is attached to.
	SetEntityID(entityID string) MediaInterface

	// GetTitle returns the media title.
	GetTitle() string
	// SetTitle sets the media title.
	SetTitle(title string) MediaInterface

	// GetDescription returns the media description.
	GetDescription() string
	// SetDescription sets the media description.
	SetDescription(description string) MediaInterface

	// GetMemo returns the internal memo.
	GetMemo() string
	// SetMemo sets the internal memo.
	SetMemo(memo string) MediaInterface

	// GetURL returns the file URL.
	GetURL() string
	// SetURL sets the file URL.
	SetURL(url string) MediaInterface

	// GetType returns the file type (mime type).
	GetType() string
	// SetType sets the file type (mime type).
	SetType(fileType string) MediaInterface

	// GetSize returns the file size as a string.
	GetSize() string
	// SetSize sets the file size as a string.
	SetSize(size string) MediaInterface

	// GetExtension returns the file extension.
	GetExtension() string
	// SetExtension sets the file extension.
	SetExtension(extension string) MediaInterface

	// GetSequence returns the display sequence/order of the file.
	GetSequence() int
	// SetSequence sets the display sequence/order of the file.
	SetSequence(sequence int) MediaInterface

	// GetStatus returns the current status.
	GetStatus() string
	// SetStatus sets the current status.
	SetStatus(status string) MediaInterface

	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) MediaInterface

	// GetUpdatedAt returns the last update timestamp as a string.
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) MediaInterface

	// GetSoftDeletedAt returns the soft deletion timestamp as a string.
	GetSoftDeletedAt() string
	// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a carbon.Carbon instance.
	GetSoftDeletedAtCarbon() *carbon.Carbon
	// SetSoftDeletedAt sets the soft deletion timestamp.
	SetSoftDeletedAt(softDeletedAt string) MediaInterface

	// Metadata methods

	// GetMetas returns all metadata as a map.
	GetMetas() (map[string]string, error)
	// GetMeta returns a specific metadata value by name.
	GetMeta(name string) string
	// SetMeta sets a single metadata value.
	SetMeta(name string, value string) error
	// SetMetas replaces all metadata with the provided map.
	SetMetas(metas map[string]string) error
	// MetasUpsert merges the provided metadata with existing values.
	MetasUpsert(metas map[string]string) error
	// MetaRemove removes a single metadata entry.
	MetaRemove(name string) error
	// MetasRemove removes multiple metadata entries.
	MetasRemove(names []string) error

	// Status predicates

	// IsActive returns true if the media status is active.
	IsActive() bool
	// IsDraft returns true if the media status is draft.
	IsDraft() bool
	// IsInactive returns true if the media status is inactive.
	IsInactive() bool

	// Type predicates

	// IsImage returns true if the media type starts with "image/".
	IsImage() bool
	// IsVideo returns true if the media type starts with "video/".
	IsVideo() bool

	// GetData returns all media data as a map.
	GetData() map[string]string
}

// Compile-time check to ensure mediaImplementation implements MediaInterface.
var _ MediaInterface = (*mediaImplementation)(nil)

// NewMedia creates a new Media instance with default values.
// The media is initialized with a generated ID, empty fields, current timestamps,
// draft status, and the max datetime for soft deletion (not deleted).
func NewMedia() MediaInterface {
	o := &mediaImplementation{}
	o.SetID(GenerateShortID()).
		SetEntityID("").
		SetTitle("").
		SetDescription("").
		SetMemo("").
		SetURL("").
		SetType("").
		SetSize("0").
		SetExtension("").
		SetSequence(0).
		SetStatus(MEDIA_STATUS_DRAFT).
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(MAX_DATETIME)

	_ = o.SetMetas(map[string]string{})

	return o
}

// mediaImplementation is the concrete implementation of MediaInterface.
// It uses neat ORM traits for data storage.
type mediaImplementation struct {
	orm.ShortID
	orm.CreatedAt
	orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate

	EntityID    string `db:"entity_id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Memo        string `db:"memo"`
	URL         string `db:"media_url"`
	Type        string `db:"media_type"`
	Size        string `db:"file_size"`
	Extension   string `db:"file_extension"`
	Sequence    int    `db:"sequence"`
	Status      string `db:"status"`
	Metas       string `db:"metas"`
}

// IsSoftDeleted returns true if the media is soft deleted.
func (o *mediaImplementation) IsSoftDeleted() bool {
	return o.SoftDeletesMaxDate.IsSoftDeleted()
}

// GetID returns the unique identifier of the media.
func (o *mediaImplementation) GetID() string {
	return o.ShortID.ID
}

// SetID sets the unique identifier of the media.
func (o *mediaImplementation) SetID(id string) MediaInterface {
	o.ShortID.ID = id
	return o
}

// GetEntityID returns the ID of the entity this media is attached to.
func (o *mediaImplementation) GetEntityID() string {
	return o.EntityID
}

// SetEntityID sets the ID of the entity this media is attached to.
func (o *mediaImplementation) SetEntityID(entityID string) MediaInterface {
	o.EntityID = entityID
	return o
}

// GetTitle returns the media title.
func (o *mediaImplementation) GetTitle() string {
	return o.Title
}

// SetTitle sets the media title.
func (o *mediaImplementation) SetTitle(title string) MediaInterface {
	o.Title = title
	return o
}

// GetDescription returns the media description.
func (o *mediaImplementation) GetDescription() string {
	return o.Description
}

// SetDescription sets the media description.
func (o *mediaImplementation) SetDescription(description string) MediaInterface {
	o.Description = description
	return o
}

// GetMemo returns the internal memo.
func (o *mediaImplementation) GetMemo() string {
	return o.Memo
}

// SetMemo sets the internal memo.
func (o *mediaImplementation) SetMemo(memo string) MediaInterface {
	o.Memo = memo
	return o
}

// GetURL returns the file URL.
func (o *mediaImplementation) GetURL() string {
	return o.URL
}

// SetURL sets the file URL.
func (o *mediaImplementation) SetURL(url string) MediaInterface {
	o.URL = url
	return o
}

// GetType returns the file type (mime type).
func (o *mediaImplementation) GetType() string {
	return o.Type
}

// SetType sets the file type (mime type).
func (o *mediaImplementation) SetType(fileType string) MediaInterface {
	o.Type = fileType
	return o
}

// GetSize returns the file size as a string.
func (o *mediaImplementation) GetSize() string {
	return o.Size
}

// SetSize sets the file size as a string.
func (o *mediaImplementation) SetSize(size string) MediaInterface {
	o.Size = size
	return o
}

// GetExtension returns the file extension.
func (o *mediaImplementation) GetExtension() string {
	return o.Extension
}

// SetExtension sets the file extension.
func (o *mediaImplementation) SetExtension(extension string) MediaInterface {
	o.Extension = extension
	return o
}

// GetSequence returns the display sequence/order of the file.
func (o *mediaImplementation) GetSequence() int {
	return o.Sequence
}

// SetSequence sets the display sequence/order of the file.
func (o *mediaImplementation) SetSequence(sequence int) MediaInterface {
	o.Sequence = sequence
	return o
}

// GetStatus returns the current status.
func (o *mediaImplementation) GetStatus() string {
	return o.Status
}

// SetStatus sets the current status.
func (o *mediaImplementation) SetStatus(status string) MediaInterface {
	o.Status = status
	return o
}

// GetCreatedAt returns the creation timestamp as a string.
func (o *mediaImplementation) GetCreatedAt() string {
	if o.CreatedAt.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAt.CreatedAt).ToDateTimeString()
}

// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
func (o *mediaImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAt.CreatedAt)
}

// SetCreatedAt sets the creation timestamp.
func (o *mediaImplementation) SetCreatedAt(createdAt string) MediaInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAt.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// GetUpdatedAt returns the last update timestamp as a string.
func (o *mediaImplementation) GetUpdatedAt() string {
	if o.UpdatedAt.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAt.UpdatedAt).ToDateTimeString()
}

// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
func (o *mediaImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAt.UpdatedAt)
}

// SetUpdatedAt sets the last update timestamp.
func (o *mediaImplementation) SetUpdatedAt(updatedAt string) MediaInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAt.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

// GetSoftDeletedAt returns the soft deletion timestamp as a string.
func (o *mediaImplementation) GetSoftDeletedAt() string {
	if o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a carbon.Carbon instance.
func (o *mediaImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (o *mediaImplementation) SetSoftDeletedAt(softDeletedAt string) MediaInterface {
	if softDeletedAt == "" {
		return o
	}
	o.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	return o
}

// GetMetas returns all metadata as a map. Returns empty map if no metas stored.
func (o *mediaImplementation) GetMetas() (map[string]string, error) {
	metasStr := o.Metas

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson := map[string]string{}
	errJson := json.Unmarshal([]byte(metasStr), &metasJson)
	if errJson != nil {
		return map[string]string{}, errJson
	}

	if metasJson == nil {
		metasJson = map[string]string{}
	}

	return metasJson, nil
}

// GetMeta returns a specific metadata value by name. Returns empty string if not found.
func (o *mediaImplementation) GetMeta(name string) string {
	metas, err := o.GetMetas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

// SetMeta sets a single metadata value.
func (o *mediaImplementation) SetMeta(name string, value string) error {
	return o.MetasUpsert(map[string]string{name: value})
}

// SetMetas replaces all metadata with the provided map.
func (o *mediaImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	o.Metas = string(mapString)
	return nil
}

// MetasUpsert merges the provided metadata with existing values.
func (o *mediaImplementation) MetasUpsert(metas map[string]string) error {
	currentMetas, err := o.GetMetas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

// MetaRemove removes a single metadata entry.
func (o *mediaImplementation) MetaRemove(name string) error {
	metas, err := o.GetMetas()
	if err != nil {
		return err
	}
	delete(metas, name)
	return o.SetMetas(metas)
}

// MetasRemove removes multiple metadata entries.
func (o *mediaImplementation) MetasRemove(names []string) error {
	for _, name := range names {
		if err := o.MetaRemove(name); err != nil {
			return err
		}
	}
	return nil
}

// IsActive returns true if the media status is active.
func (o *mediaImplementation) IsActive() bool {
	return o.Status == MEDIA_STATUS_ACTIVE
}

// IsDraft returns true if the media status is draft.
func (o *mediaImplementation) IsDraft() bool {
	return o.Status == MEDIA_STATUS_DRAFT
}

// IsInactive returns true if the media status is inactive.
func (o *mediaImplementation) IsInactive() bool {
	return o.Status == MEDIA_STATUS_INACTIVE
}

// IsImage returns true if the media type starts with "image/".
func (o *mediaImplementation) IsImage() bool {
	return strings.HasPrefix(o.Type, "image/")
}

// IsVideo returns true if the media type starts with "video/".
func (o *mediaImplementation) IsVideo() bool {
	return strings.HasPrefix(o.Type, "video/")
}

// GetData returns all media data as a map.
func (o *mediaImplementation) GetData() map[string]string {
	var createdAt, updatedAt, softDeletedAt string
	if !o.CreatedAt.CreatedAt.IsZero() {
		createdAt = carbon.CreateFromStdTime(o.CreatedAt.CreatedAt).ToDateTimeString(carbon.UTC)
	}
	if !o.UpdatedAt.UpdatedAt.IsZero() {
		updatedAt = carbon.CreateFromStdTime(o.UpdatedAt.UpdatedAt).ToDateTimeString(carbon.UTC)
	}
	if !o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		softDeletedAt = carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString(carbon.UTC)
	}

	return map[string]string{
		COLUMN_ID:              o.ShortID.ID,
		COLUMN_ENTITY_ID:       o.EntityID,
		COLUMN_TITLE:           o.Title,
		COLUMN_DESCRIPTION:     o.Description,
		COLUMN_MEMO:            o.Memo,
		COLUMN_MEDIA_URL:       o.URL,
		COLUMN_MEDIA_TYPE:      o.Type,
		COLUMN_FILE_SIZE:       o.Size,
		COLUMN_FILE_EXTENSION:  o.Extension,
		COLUMN_SEQUENCE:        strconv.Itoa(o.Sequence),
		COLUMN_STATUS:          o.Status,
		COLUMN_METAS:           o.Metas,
		COLUMN_CREATED_AT:      createdAt,
		COLUMN_UPDATED_AT:      updatedAt,
		COLUMN_SOFT_DELETED_AT: softDeletedAt,
	}
}
