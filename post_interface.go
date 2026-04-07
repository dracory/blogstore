package blogstore

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// PostInterface defines the interface for blog post operations
type PostInterface interface {
	// Identity
	ID() string
	SetID(id string) PostInterface

	// Author
	AuthorID() string
	SetAuthorID(authorID string) PostInterface

	// Content
	Title() string
	SetTitle(title string) PostInterface
	Slug() string

	Content() string
	SetContent(content string) PostInterface

	Summary() string
	SetSummary(summary string) PostInterface

	// Content Type and Editor
	ContentType() string
	SetContentType(contentType string) PostInterface
	Editor() string
	SetEditor(editor string) PostInterface

	IsContentMarkdown() bool
	IsContentHtml() bool
	IsContentPlainText() bool
	IsContentBlocks() bool

	// SEO and Meta
	CanonicalURL() string
	SetCanonicalURL(canonicalURL string) PostInterface

	MetaDescription() string
	SetMetaDescription(metaDescription string) PostInterface

	MetaKeywords() string
	SetMetaKeywords(metaKeywords string) PostInterface

	MetaRobots() string
	SetMetaRobots(metaRobots string) PostInterface

	// Featured Image
	ImageUrl() string
	SetImageUrl(imageURL string) PostInterface
	ImageUrlOrDefault() string

	// Status
	Status() string
	SetStatus(status string) PostInterface

	IsDraft() bool
	IsPublished() bool
	IsUnpublished() bool
	IsTrashed() bool

	// Publishing
	PublishedAt() string
	SetPublishedAt(publishedAt string) PostInterface
	PublishedAtCarbon() *carbon.Carbon
	PublishedAtTime() time.Time

	// Timestamps
	CreatedAt() string
	SetCreatedAt(createdAt string) PostInterface
	CreatedAtCarbon() *carbon.Carbon
	CreatedAtTime() time.Time

	UpdatedAt() string
	SetUpdatedAt(updatedAt string) PostInterface
	UpdatedAtCarbon() *carbon.Carbon

	SoftDeletedAt() string
	SetSoftDeletedAt(deletedAt string) PostInterface
	SoftDeletedAtCarbon() *carbon.Carbon

	// Memo
	Memo() string
	SetMemo(memo string) PostInterface

	// Featured flag
	Featured() string
	SetFeatured(featured string) PostInterface

	// Metadata
	Meta(key string) string
	SetMeta(key string, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	AddMetas(metas map[string]string) error

	// Versioning
	MarshalToVersioning() (string, error)
	UnmarshalFromVersioning(content string) error

	// DataObject methods (from embedded dataobject.DataObject)
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()
	Get(key string) string
	Set(key string, value string)
	Hydrate(data map[string]string)
	IsDirty() bool
}

// Ensure Post implements PostInterface
var _ PostInterface = (*Post)(nil)
