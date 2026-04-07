package blogstore

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// PostInterface defines the interface for blog post operations
type PostInterface interface {
	// Identity
	GetID() string
	SetID(id string) PostInterface

	// Author
	GetAuthorID() string
	SetAuthorID(authorID string) PostInterface

	// Content
	GetTitle() string
	SetTitle(title string) PostInterface
	GetSlug() string

	GetContent() string
	SetContent(content string) PostInterface

	GetSummary() string
	SetSummary(summary string) PostInterface

	// Content Type and Editor
	GetContentType() string
	SetContentType(contentType string) PostInterface
	GetEditor() string
	SetEditor(editor string) PostInterface

	IsContentMarkdown() bool
	IsContentHtml() bool
	IsContentPlainText() bool
	IsContentBlocks() bool

	// SEO and Meta
	GetCanonicalURL() string
	SetCanonicalURL(canonicalURL string) PostInterface

	GetMetaDescription() string
	SetMetaDescription(metaDescription string) PostInterface

	GetMetaKeywords() string
	SetMetaKeywords(metaKeywords string) PostInterface

	GetMetaRobots() string
	SetMetaRobots(metaRobots string) PostInterface

	// Featured Image
	GetImageUrl() string
	SetImageUrl(imageURL string) PostInterface
	GetImageUrlOrDefault() string

	// Status
	GetStatus() string
	SetStatus(status string) PostInterface

	IsDraft() bool
	IsPublished() bool
	IsUnpublished() bool
	IsTrashed() bool

	// Publishing
	GetPublishedAt() string
	SetPublishedAt(publishedAt string) PostInterface
	GetPublishedAtCarbon() *carbon.Carbon
	GetPublishedAtTime() time.Time

	// Timestamps
	GetCreatedAt() string
	SetCreatedAt(createdAt string) PostInterface
	GetCreatedAtCarbon() *carbon.Carbon
	GetCreatedAtTime() time.Time

	GetUpdatedAt() string
	SetUpdatedAt(updatedAt string) PostInterface
	GetUpdatedAtCarbon() *carbon.Carbon

	GetSoftDeletedAt() string
	SetSoftDeletedAt(deletedAt string) PostInterface
	GetSoftDeletedAtCarbon() *carbon.Carbon

	// Memo
	GetMemo() string
	SetMemo(memo string) PostInterface

	// Featured flag
	GetFeatured() string
	SetFeatured(featured string) PostInterface

	// Metadata
	GetMeta(key string) string
	SetMeta(key string, value string) error
	GetMetas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	AddMetas(metas map[string]string) error

	// Versioning
	MarshalToVersioning() (string, error)
	UnmarshalFromVersioning(content string) error

	// DataObject methods (from embedded dataobject.DataObject)
	GetData() map[string]string
	GetDataChanged() map[string]string
	MarkAsNotDirty()
	Get(key string) string
	Set(key string, value string)
	Hydrate(data map[string]string)
	IsDirty() bool
}

// Ensure postImplementation implements PostInterface
var _ PostInterface = (*postImplementation)(nil)
