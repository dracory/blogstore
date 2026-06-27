package blogstore

import (
	"strconv"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	"github.com/dromara/carbon/v2"
)

// PostFileInterface defines the interface for post file (media attachment) operations.
// Post files represent media files attached to a post, such as images, videos, or documents.
type PostFileInterface interface {
	// IsSoftDeleted returns true if the post file is soft deleted.
	IsSoftDeleted() bool

	// GetID returns the unique identifier of the post file.
	GetID() string
	// SetID sets the unique identifier of the post file.
	SetID(id string) PostFileInterface

	// GetPostID returns the ID of the post this file is attached to.
	GetPostID() string
	// SetPostID sets the ID of the post this file is attached to.
	SetPostID(postID string) PostFileInterface

	// GetName returns the file name.
	GetName() string
	// SetName sets the file name.
	SetName(name string) PostFileInterface

	// GetURL returns the file URL.
	GetURL() string
	// SetURL sets the file URL.
	SetURL(url string) PostFileInterface

	// GetType returns the file type (mime type).
	GetType() string
	// SetType sets the file type (mime type).
	SetType(fileType string) PostFileInterface

	// GetSize returns the file size as a string.
	GetSize() string
	// SetSize sets the file size as a string.
	SetSize(size string) PostFileInterface

	// GetExtension returns the file extension.
	GetExtension() string
	// SetExtension sets the file extension.
	SetExtension(extension string) PostFileInterface

	// GetSequence returns the display sequence/order of the file.
	GetSequence() int
	// SetSequence sets the display sequence/order of the file.
	SetSequence(sequence int) PostFileInterface

	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) PostFileInterface

	// GetUpdatedAt returns the last update timestamp as a string.
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) PostFileInterface

	// GetSoftDeletedAt returns the soft deletion timestamp as a string.
	GetSoftDeletedAt() string
	// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a carbon.Carbon instance.
	GetSoftDeletedAtCarbon() *carbon.Carbon
	// SetSoftDeletedAt sets the soft deletion timestamp.
	SetSoftDeletedAt(softDeletedAt string) PostFileInterface

	// GetData returns all post file data as a map.
	GetData() map[string]string
}

// Compile-time check to ensure postFileImplementation implements PostFileInterface.
var _ PostFileInterface = (*postFileImplementation)(nil)

// NewPostFile creates a new PostFile instance with default values.
// The post file is initialized with a generated ID, empty fields, current timestamps,
// and the max datetime for soft deletion (not deleted).
func NewPostFile() PostFileInterface {
	o := &postFileImplementation{}
	o.SetID(GenerateShortID()).
		SetPostID("").
		SetName("").
		SetURL("").
		SetType("").
		SetSize("0").
		SetExtension("").
		SetSequence(0).
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(MAX_DATETIME)

	return o
}

// NewPostFileFromExistingData creates a PostFile instance from existing data.
// This is useful when hydrating a post file from database records.
func NewPostFileFromExistingData(data map[string]string) PostFileInterface {
	o := &postFileImplementation{}
	if v, ok := data[COLUMN_ID]; ok {
		o.SetID(v)
	}
	if v, ok := data[COLUMN_POST_ID]; ok {
		o.SetPostID(v)
	}
	if v, ok := data[COLUMN_NAME]; ok {
		o.SetName(v)
	}
	if v, ok := data[COLUMN_URL]; ok {
		o.SetURL(v)
	}
	if v, ok := data[COLUMN_FILE_TYPE]; ok {
		o.SetType(v)
	}
	if v, ok := data[COLUMN_FILE_SIZE]; ok {
		o.SetSize(v)
	}
	if v, ok := data[COLUMN_FILE_EXTENSION]; ok {
		o.SetExtension(v)
	}
	if v, ok := data[COLUMN_SEQUENCE]; ok {
		if seq, err := strconv.Atoi(v); err == nil {
			o.SetSequence(seq)
		}
	}
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		o.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		o.SetUpdatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		o.SetSoftDeletedAt(v)
	}
	return o
}

// postFileImplementation is the concrete implementation of PostFileInterface.
// It uses neat ORM traits for data storage.
type postFileImplementation struct {
	orm.ShortID
	soft_delete.SoftDeletesMaxDate

	PostIDField    string `db:"post_id"`
	NameField      string `db:"name"`
	URLField       string `db:"url"`
	TypeField      string `db:"file_type"`
	SizeField      string `db:"file_size"`
	ExtensionField string `db:"file_extension"`
	SequenceField  int    `db:"sequence"`
	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
}

// IsSoftDeleted returns true if the post file is soft deleted.
func (o *postFileImplementation) IsSoftDeleted() bool {
	return o.SoftDeletesMaxDate.IsSoftDeleted()
}

// GetID returns the unique identifier of the post file.
func (o *postFileImplementation) GetID() string {
	return o.ShortID.ID
}

// SetID sets the unique identifier of the post file.
func (o *postFileImplementation) SetID(id string) PostFileInterface {
	o.ShortID.ID = id
	return o
}

// GetPostID returns the ID of the post this file is attached to.
func (o *postFileImplementation) GetPostID() string {
	return o.PostIDField
}

// SetPostID sets the ID of the post this file is attached to.
func (o *postFileImplementation) SetPostID(postID string) PostFileInterface {
	o.PostIDField = postID
	return o
}

// GetName returns the file name.
func (o *postFileImplementation) GetName() string {
	return o.NameField
}

// SetName sets the file name.
func (o *postFileImplementation) SetName(name string) PostFileInterface {
	o.NameField = name
	return o
}

// GetURL returns the file URL.
func (o *postFileImplementation) GetURL() string {
	return o.URLField
}

// SetURL sets the file URL.
func (o *postFileImplementation) SetURL(url string) PostFileInterface {
	o.URLField = url
	return o
}

// GetType returns the file type (mime type).
func (o *postFileImplementation) GetType() string {
	return o.TypeField
}

// SetType sets the file type (mime type).
func (o *postFileImplementation) SetType(fileType string) PostFileInterface {
	o.TypeField = fileType
	return o
}

// GetSize returns the file size as a string.
func (o *postFileImplementation) GetSize() string {
	return o.SizeField
}

// SetSize sets the file size as a string.
func (o *postFileImplementation) SetSize(size string) PostFileInterface {
	o.SizeField = size
	return o
}

// GetExtension returns the file extension.
func (o *postFileImplementation) GetExtension() string {
	return o.ExtensionField
}

// SetExtension sets the file extension.
func (o *postFileImplementation) SetExtension(extension string) PostFileInterface {
	o.ExtensionField = extension
	return o
}

// GetSequence returns the display sequence/order of the file.
func (o *postFileImplementation) GetSequence() int {
	return o.SequenceField
}

// SetSequence sets the display sequence/order of the file.
func (o *postFileImplementation) SetSequence(sequence int) PostFileInterface {
	o.SequenceField = sequence
	return o
}

// GetCreatedAt returns the creation timestamp as a string.
func (o *postFileImplementation) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
func (o *postFileImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

// SetCreatedAt sets the creation timestamp.
func (o *postFileImplementation) SetCreatedAt(createdAt string) PostFileInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// GetUpdatedAt returns the last update timestamp as a string.
func (o *postFileImplementation) GetUpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
func (o *postFileImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

// SetUpdatedAt sets the last update timestamp.
func (o *postFileImplementation) SetUpdatedAt(updatedAt string) PostFileInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

// GetSoftDeletedAt returns the soft deletion timestamp as a string.
func (o *postFileImplementation) GetSoftDeletedAt() string {
	if o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a carbon.Carbon instance.
func (o *postFileImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (o *postFileImplementation) SetSoftDeletedAt(softDeletedAt string) PostFileInterface {
	if softDeletedAt == "" {
		return o
	}
	o.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	return o
}

// GetData returns all post file data as a map.
func (o *postFileImplementation) GetData() map[string]string {
	var createdAt, updatedAt, softDeletedAt string
	if !o.CreatedAtField.CreatedAt.IsZero() {
		createdAt = carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString(carbon.UTC)
	}
	if !o.UpdatedAtField.UpdatedAt.IsZero() {
		updatedAt = carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString(carbon.UTC)
	}
	if !o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		softDeletedAt = carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString(carbon.UTC)
	}

	return map[string]string{
		COLUMN_ID:              o.ShortID.ID,
		COLUMN_POST_ID:         o.PostIDField,
		COLUMN_NAME:            o.NameField,
		COLUMN_URL:             o.URLField,
		COLUMN_FILE_TYPE:       o.TypeField,
		COLUMN_FILE_SIZE:       o.SizeField,
		COLUMN_FILE_EXTENSION:  o.ExtensionField,
		COLUMN_SEQUENCE:        strconv.Itoa(o.SequenceField),
		COLUMN_CREATED_AT:      createdAt,
		COLUMN_UPDATED_AT:      updatedAt,
		COLUMN_SOFT_DELETED_AT: softDeletedAt,
	}
}
