package blogstore

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// PostInterface defines the interface for blog post operations.
// Posts represent the main content entity in the blog system with support for
// multiple content types, versioning, taxonomies, and SEO metadata.
type PostInterface interface {
	// Identity
	// GetID returns the unique identifier of the post.
	GetID() string
	// SetID sets the unique identifier of the post.
	SetID(id string) PostInterface

	// Author
	// GetAuthorID returns the ID of the post author.
	GetAuthorID() string
	// SetAuthorID sets the ID of the post author.
	SetAuthorID(authorID string) PostInterface

	// Content
	// GetTitle returns the post title.
	GetTitle() string
	// SetTitle sets the post title.
	SetTitle(title string) PostInterface
	// GetSlug returns the URL-friendly slug generated from the title.
	GetSlug() string

	// GetContent returns the main content/body of the post.
	GetContent() string
	// SetContent sets the main content/body of the post.
	SetContent(content string) PostInterface

	// GetSummary returns the post summary/excerpt.
	GetSummary() string
	// SetSummary sets the post summary/excerpt.
	SetSummary(summary string) PostInterface

	// Content Type and Editor
	// GetContentType returns the content type of this post (markdown, html, plain_text, blocks).
	GetContentType() string
	// SetContentType sets the content type of this post.
	SetContentType(contentType string) PostInterface
	// GetEditor returns the editor type for this post.
	GetEditor() string
	// SetEditor sets the editor type for this post.
	SetEditor(editor string) PostInterface

	// IsContentMarkdown returns true if the post content type is markdown.
	IsContentMarkdown() bool
	// IsContentHtml returns true if the post content type is HTML.
	IsContentHtml() bool
	// IsContentPlainText returns true if the post content type is plain text.
	IsContentPlainText() bool
	// IsContentBlocks returns true if the post content type is blocks.
	IsContentBlocks() bool

	// SEO and Meta
	// GetCanonicalURL returns the canonical URL for SEO purposes.
	GetCanonicalURL() string
	// SetCanonicalURL sets the canonical URL for SEO purposes.
	SetCanonicalURL(canonicalURL string) PostInterface

	// GetMetaDescription returns the SEO meta description.
	GetMetaDescription() string
	// SetMetaDescription sets the SEO meta description.
	SetMetaDescription(metaDescription string) PostInterface

	// GetMetaKeywords returns the SEO meta keywords.
	GetMetaKeywords() string
	// SetMetaKeywords sets the SEO meta keywords.
	SetMetaKeywords(metaKeywords string) PostInterface

	// GetMetaRobots returns the SEO robots meta tag value.
	GetMetaRobots() string
	// SetMetaRobots sets the SEO robots meta tag value.
	SetMetaRobots(metaRobots string) PostInterface

	// Featured Image
	// GetImageUrl returns the URL of the post's featured image.
	GetImageUrl() string
	// SetImageUrl sets the URL of the post's featured image.
	SetImageUrl(imageURL string) PostInterface
	// GetImageUrlOrDefault returns the featured image URL, or a default URL if none is set.
	GetImageUrlOrDefault() string

	// Status
	// GetStatus returns the post status (draft, published, trash, etc.).
	GetStatus() string
	// SetStatus sets the post status (draft, published, trash, etc.).
	SetStatus(status string) PostInterface

	// IsDraft returns true if the post status is POST_STATUS_DRAFT.
	IsDraft() bool
	// IsPublished returns true if the post status is POST_STATUS_PUBLISHED.
	IsPublished() bool
	// IsUnpublished returns true if the post status is not published.
	IsUnpublished() bool
	// IsTrashed returns true if the post status is POST_STATUS_TRASH.
	IsTrashed() bool

	// Publishing
	// GetPublishedAt returns the publication timestamp as a string.
	GetPublishedAt() string
	// SetPublishedAt sets the publication timestamp.
	SetPublishedAt(publishedAt string) PostInterface
	// GetPublishedAtCarbon returns the publication timestamp as a carbon.Carbon instance.
	GetPublishedAtCarbon() *carbon.Carbon
	// GetPublishedAtTime returns the publication timestamp as a time.Time instance.
	GetPublishedAtTime() time.Time

	// Timestamps
	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) PostInterface
	// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
	GetCreatedAtTime() time.Time

	// GetUpdatedAt returns the last update timestamp as a string.
	GetUpdatedAt() string
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) PostInterface
	// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon

	// GetSoftDeletedAt returns the soft deletion timestamp as a string.
	GetSoftDeletedAt() string
	// SetSoftDeletedAt sets the soft deletion timestamp.
	SetSoftDeletedAt(deletedAt string) PostInterface
	// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a carbon.Carbon instance.
	GetSoftDeletedAtCarbon() *carbon.Carbon

	// Memo
	// GetMemo returns the internal memo/note for the post.
	GetMemo() string
	// SetMemo sets the internal memo/note for the post.
	SetMemo(memo string) PostInterface

	// Featured flag
	// GetFeatured returns the featured status (YES/NO) of the post.
	GetFeatured() string
	// SetFeatured sets the featured status (YES/NO) of the post.
	SetFeatured(featured string) PostInterface

	// Metadata
	// GetMeta retrieves a single metadata value by key.
	GetMeta(key string) string
	// SetMeta sets a single metadata value by key.
	SetMeta(key string, value string) error
	// GetMetas returns all metadata as a map[string]string.
	GetMetas() (map[string]string, error)
	// SetMetas sets all metadata from a map[string]string.
	SetMetas(metas map[string]string) error
	// AddMetas adds multiple metadata key-value pairs to the existing metas.
	AddMetas(metas map[string]string) error

	// Versioning
	// MarshalToVersioning serializes the post data for versioning storage.
	MarshalToVersioning() (string, error)
	// UnmarshalFromVersioning restores post data from a serialized versioning string.
	UnmarshalFromVersioning(content string) error

	// Taxonomy methods
	// TermIDs retrieves the term IDs for a specific taxonomy from the post metadata.
	TermIDs(taxonomySlug string) []string
	// SetTermIDs stores the term IDs for a specific taxonomy in the post metadata.
	SetTermIDs(taxonomySlug string, termIDs []string) PostInterface

	// Convenience helpers
	// CategoryIDs retrieves the category IDs associated with this post.
	CategoryIDs() []string
	// SetCategoryIDs sets the category IDs for this post.
	SetCategoryIDs(ids []string) PostInterface
	// TagIDs retrieves the tag IDs associated with this post.
	TagIDs() []string
	// SetTagIDs sets the tag IDs for this post.
	SetTagIDs(ids []string) PostInterface

	// DataObject methods (from embedded dataobject.DataObject)
	// GetData returns all post data as a map.
	GetData() map[string]string
	// GetDataChanged returns only the fields that have been modified.
	GetDataChanged() map[string]string
	// MarkAsNotDirty clears the dirty state of the post.
	MarkAsNotDirty()
	// Get retrieves a value by key from the post data.
	Get(key string) string
	// Set stores a value by key in the post data.
	Set(key string, value string)
	// Hydrate populates the post with data from a map.
	Hydrate(data map[string]string)
	// IsDirty returns true if the post has unsaved changes.
	IsDirty() bool
}

// Compile-time check to ensure postImplementation implements PostInterface.
var _ PostInterface = (*postImplementation)(nil)
