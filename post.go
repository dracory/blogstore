package blogstore

import (
	"encoding/json"
	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	"github.com/dracory/str"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
	"time"
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
	// SetSlug sets the URL-friendly slug for this post.
	SetSlug(slug string) PostInterface

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

	// Old Slug methods (WordPress-style slug history for redirects)
	// GetOldSlugs retrieves the array of historical slugs for redirect purposes.
	GetOldSlugs() []string
	// SetOldSlugs sets the array of historical slugs.
	SetOldSlugs(slugs []string) error
	// AddOldSlug adds a slug to the old slugs history.
	AddOldSlug(slug string) error

	// DataObject methods (from embedded dataobject.DataObject)
	// GetData returns all post data as a map.
	GetData() map[string]string
	// GetDataChanged returns only the fields that have been modified.
	GetDataChanged() map[string]string
	// MarkAsNotDirty clears the dirty state of the post.
	// If no columns specified, marks all fields as not dirty.
	// If columns specified, marks only those columns as not dirty.
	MarkAsNotDirty(columns ...string)
	// MarkAsDirty marks the post as dirty.
	// If no columns specified, marks all fields as dirty.
	// If columns specified, marks only those columns as dirty.
	MarkAsDirty(columns ...string)
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

// NewPost creates a new Post instance with default values.
// The post is initialized with a generated ID, draft status, empty content fields,
// current timestamps, and an empty metadata map.
func NewPost() PostInterface {
	o := &postImplementation{}
	o.SetID(GenerateShortID()).
		SetAuthorID("").
		SetCanonicalURL("").
		SetContent("").
		SetFeatured(NO).
		SetImageUrl("").
		SetMemo("").
		SetMetaDescription("").
		SetMetaKeywords("").
		SetMetaRobots("").
		SetSlug("").
		SetStatus(POST_STATUS_DRAFT).
		SetSummary("").
		SetTitle("").
		SetPublishedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(MAX_DATETIME).
		SetMetas(map[string]string{})

	return o
}

// NewPostFromExistingData creates a Post instance from existing data.
// This is useful when hydrating a post from database records.
func NewPostFromExistingData(data map[string]string) PostInterface {
	o := &postImplementation{}
	o.Hydrate(data)
	return o
}

// postImplementation is the concrete implementation of the PostInterface.
// It uses neat ORM traits for data storage.
type postImplementation struct {
	orm.ShortID

	AuthorIDField        string    `db:"author_id"`
	CanonicalURLField    string    `db:"canonical_url"`
	ContentField         string    `db:"content"`
	FeaturedField        string    `db:"featured"`
	ImageURLField        string    `db:"image_url"`
	MemoField            string    `db:"memo"`
	MetaDescriptionField string    `db:"meta_description"`
	MetaKeywordsField    string    `db:"meta_keywords"`
	MetaRobotsField      string    `db:"meta_robots"`
	MetasField           string    `db:"metas"`
	PublishedAtField     time.Time `db:"published_at"`
	StatusField          string    `db:"status"`
	SlugField            string    `db:"slug"`
	SummaryField         string    `db:"summary"`
	TitleField           string    `db:"title"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate
}

// ================================== METHODS ==================================

// GetSlug returns the URL-friendly slug for this post.
// If a custom slug is set, it returns that; otherwise, it generates one from the title.
func (o *postImplementation) GetSlug() string {
	storedSlug := o.Get(COLUMN_SLUG)
	if storedSlug != "" {
		return storedSlug
	}
	return str.Slugify(o.GetTitle(), '-')
}

// SetSlug sets the URL-friendly slug for this post.
func (o *postImplementation) SetSlug(slug string) PostInterface {
	o.Set(COLUMN_SLUG, slug)
	return o
}

// GetEditor returns the editor type for this post (e.g., markdown, html, blocks).
func (o *postImplementation) GetEditor() string {
	return o.GetMeta("editor")
}

// SetEditor sets the editor type for this post.
func (o *postImplementation) SetEditor(editor string) PostInterface {
	o.SetMeta("editor", editor)
	return o
}

// GetContentType returns the content type of this post (markdown, html, plain_text, blocks).
func (o *postImplementation) GetContentType() string {
	return o.GetMeta("content_type")
}

// SetContentType sets the content type of this post.
func (o *postImplementation) SetContentType(contentType string) PostInterface {
	o.SetMeta("content_type", contentType)
	return o
}

// IsDraft returns true if the post status is POST_STATUS_DRAFT.
func (o *postImplementation) IsDraft() bool {
	return o.GetStatus() == POST_STATUS_DRAFT
}

// IsPublished returns true if the post status is POST_STATUS_PUBLISHED.
func (o *postImplementation) IsPublished() bool {
	return o.GetStatus() == POST_STATUS_PUBLISHED
}

// IsContentMarkdown returns true if the post content type is markdown.
func (o *postImplementation) IsContentMarkdown() bool {
	return o.GetContentType() == POST_CONTENT_TYPE_MARKDOWN
}

// IsContentHtml returns true if the post content type is HTML.
func (o *postImplementation) IsContentHtml() bool {
	return o.GetContentType() == POST_CONTENT_TYPE_HTML
}

// IsContentPlainText returns true if the post content type is plain text.
func (o *postImplementation) IsContentPlainText() bool {
	return o.GetContentType() == POST_CONTENT_TYPE_PLAIN_TEXT
}

// IsContentBlocks returns true if the post content type is blocks.
func (o *postImplementation) IsContentBlocks() bool {
	return o.GetContentType() == POST_CONTENT_TYPE_BLOCKS
}

// IsTrashed returns true if the post status is POST_STATUS_TRASH.
func (o *postImplementation) IsTrashed() bool {
	return o.GetStatus() == POST_STATUS_TRASH
}

// IsUnpublished returns true if the post status is not published.
func (o *postImplementation) IsUnpublished() bool {
	return !o.IsPublished()
}

// MarshalToVersioning serializes the post data for versioning storage.
// Excludes timestamp fields (created_at, updated_at, soft_deleted_at).
func (o *postImplementation) MarshalToVersioning() (string, error) {
	versionedData := map[string]string{}

	for k, v := range o.GetData() {
		if k == COLUMN_CREATED_AT ||
			k == COLUMN_UPDATED_AT ||
			k == COLUMN_SOFT_DELETED_AT {
			continue
		}
		versionedData[k] = v
	}

	b, err := json.Marshal(versionedData)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// UnmarshalFromVersioning restores post data from a serialized versioning string.
// Excludes timestamp fields and updates the updated_at timestamp.
func (o *postImplementation) UnmarshalFromVersioning(content string) error {
	versionedData := map[string]string{}
	if err := json.Unmarshal([]byte(content), &versionedData); err != nil {
		return err
	}

	for k, v := range versionedData {
		// Skip timestamp fields that shouldn't be restored from versioning
		if k == COLUMN_CREATED_AT ||
			k == COLUMN_UPDATED_AT ||
			k == COLUMN_SOFT_DELETED_AT {
			continue
		}
		o.Set(k, v)
	}

	// Update the updated_at timestamp to current time when restoring
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return nil
}

// ============================ SETTERS AND GETTERS ============================

// AddMetas adds multiple metadata key-value pairs to the existing metas.
func (o *postImplementation) AddMetas(metas map[string]string) error {
	currentMetas, err := o.GetMetas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

// GetAuthorID returns the ID of the post author.
func (o *postImplementation) GetAuthorID() string {
	return o.Get(COLUMN_AUTHOR_ID)
}

// SetAuthorID sets the ID of the post author.
func (o *postImplementation) SetAuthorID(authorID string) PostInterface {
	o.Set(COLUMN_AUTHOR_ID, authorID)
	return o
}

// GetCanonicalURL returns the canonical URL for SEO purposes.
func (o *postImplementation) GetCanonicalURL() string {
	return o.Get(COLUMN_CANONICAL_URL)
}

// SetCanonicalURL sets the canonical URL for SEO purposes.
func (o *postImplementation) SetCanonicalURL(canonicalURL string) PostInterface {
	o.Set(COLUMN_CANONICAL_URL, canonicalURL)
	return o
}

// GetContent returns the main content/body of the post.
func (o *postImplementation) GetContent() string {
	return o.Get(COLUMN_CONTENT)
}

// SetContent sets the main content/body of the post.
func (o *postImplementation) SetContent(content string) PostInterface {
	o.Set(COLUMN_CONTENT, content)
	return o
}

// GetCreatedAt returns the creation timestamp as a string.
func (o *postImplementation) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

// SetCreatedAt sets the creation timestamp.
func (o *postImplementation) SetCreatedAt(createdAt string) PostInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
// Returns the null datetime if the created_at field is empty.
func (o *postImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
// Returns zero time if the created_at field is empty.
func (o *postImplementation) GetCreatedAtTime() time.Time {
	return o.CreatedAtField.CreatedAt
}

// GetSoftDeletedAt returns the soft deletion timestamp as a string.
func (o *postImplementation) GetSoftDeletedAt() string {
	if o.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletedAt).ToDateTimeString()
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a carbon.Carbon instance.
// Returns the null datetime if the soft_deleted_at field is empty.
func (o *postImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletedAt)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (o *postImplementation) SetSoftDeletedAt(deletedAt string) PostInterface {
	if deletedAt == "" {
		o.SoftDeletedAt = soft_delete.MaxSoftDeletedAt
		return o
	}
	o.SoftDeletedAt = carbon.Parse(deletedAt, carbon.UTC).StdTime()
	return o
}

// GetFeatured returns the featured status (YES/NO) of the post.
func (o *postImplementation) GetFeatured() string {
	return o.Get(COLUMN_FEATURED)
}

// SetFeatured sets the featured status (YES/NO) of the post.
func (o *postImplementation) SetFeatured(featured string) PostInterface {
	o.Set(COLUMN_FEATURED, featured)
	return o
}

// GetID returns the unique identifier of the post.
func (o *postImplementation) GetID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the unique identifier of the post.
func (o *postImplementation) SetID(id string) PostInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// GetImageUrl returns the URL of the post's featured image.
func (o *postImplementation) GetImageUrl() string {
	return o.Get(COLUMN_IMAGE_URL)
}

// SetImageUrl sets the URL of the post's featured image.
func (o *postImplementation) SetImageUrl(imageURL string) PostInterface {
	o.Set(COLUMN_IMAGE_URL, imageURL)
	return o
}

// GetImageUrlOrDefault returns the featured image URL, or a default image URL if none is set.
func (o *postImplementation) GetImageUrlOrDefault() string {
	return lo.Ternary(o.GetImageUrl() == "", BlogNoImageUrl(), o.GetImageUrl())
}

// GetMemo returns the internal memo/note for the post.
func (o *postImplementation) GetMemo() string {
	return o.Get(COLUMN_MEMO)
}

// SetMemo sets the internal memo/note for the post.
func (o *postImplementation) SetMemo(memo string) PostInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

// GetMeta retrieves a single metadata value by key.
// Returns an empty string if the key doesn't exist or if there's an error.
func (o *postImplementation) GetMeta(key string) string {
	metas, err := o.GetMetas()

	if err != nil {
		return ""
	}

	return lo.ValueOr(metas, key, "")
}

// SetMeta sets a single metadata value by key.
func (o *postImplementation) SetMeta(key string, value string) error {
	metas, err := o.GetMetas()

	if err != nil {
		return err
	}

	metas[key] = value
	return o.SetMetas(metas)
}

// GetMetas returns all metadata as a map[string]string.
// Returns an empty map if the metadata field is empty or invalid JSON.
func (o *postImplementation) GetMetas() (map[string]string, error) {
	metasStr := o.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson := map[string]string{}
	errJson := json.Unmarshal([]byte(metasStr), &metasJson)
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return metasJson, nil
}

// SetMetas sets all metadata from a map[string]string.
func (o *postImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	o.Set(COLUMN_METAS, string(mapString))
	return nil
}

// GetMetaDescription returns the SEO meta description.
func (o *postImplementation) GetMetaDescription() string {
	return o.Get(COLUMN_META_DESCRIPTION)
}

// SetMetaDescription sets the SEO meta description.
func (o *postImplementation) SetMetaDescription(metaDescription string) PostInterface {
	o.Set(COLUMN_META_DESCRIPTION, metaDescription)
	return o
}

// GetMetaKeywords returns the SEO meta keywords.
func (o *postImplementation) GetMetaKeywords() string {
	return o.Get(COLUMN_META_KEYWORDS)
}

// SetMetaKeywords sets the SEO meta keywords.
func (o *postImplementation) SetMetaKeywords(metaKeywords string) PostInterface {
	o.Set(COLUMN_META_KEYWORDS, metaKeywords)
	return o
}

// GetMetaRobots returns the SEO robots meta tag value.
func (o *postImplementation) GetMetaRobots() string {
	return o.Get(COLUMN_META_ROBOTS)
}

// SetMetaRobots sets the SEO robots meta tag value.
func (o *postImplementation) SetMetaRobots(metaRobots string) PostInterface {
	o.Set(COLUMN_META_ROBOTS, metaRobots)
	return o
}

// GetPublishedAt returns the publication timestamp as a string.
func (o *postImplementation) GetPublishedAt() string {
	return o.Get(COLUMN_PUBLISHED_AT)
}

// SetPublishedAt sets the publication timestamp.
func (o *postImplementation) SetPublishedAt(status string) PostInterface {
	o.Set(COLUMN_PUBLISHED_AT, status)
	return o
}

// GetPublishedAtCarbon returns the publication timestamp as a carbon.Carbon instance.
// Returns the null datetime if the published_at field is empty.
func (o *postImplementation) GetPublishedAtCarbon() *carbon.Carbon {
	createdAt := o.GetPublishedAt()
	if createdAt == "" {
		return carbon.Parse(NULL_DATETIME)
	}
	return carbon.Parse(createdAt)
}

// GetPublishedAtTime returns the publication timestamp as a time.Time instance.
// Returns zero time if the published_at field is empty.
func (o *postImplementation) GetPublishedAtTime() time.Time {
	publishedAt := o.GetPublishedAt()
	if publishedAt == "" {
		return time.Time{}
	}
	return carbon.Parse(publishedAt).StdTime()
}

// GetStatus returns the post status (draft, published, trash, etc.).
func (o *postImplementation) GetStatus() string {
	return o.Get(COLUMN_STATUS)
}

// SetStatus sets the post status (draft, published, trash, etc.).
func (o *postImplementation) SetStatus(status string) PostInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

// GetSummary returns the post summary/excerpt.
func (o *postImplementation) GetSummary() string {
	return o.Get(COLUMN_SUMMARY)
}

// SetSummary sets the post summary/excerpt.
func (o *postImplementation) SetSummary(summary string) PostInterface {
	o.Set(COLUMN_SUMMARY, summary)
	return o
}

// GetTitle returns the post title.
func (o *postImplementation) GetTitle() string {
	return o.Get(COLUMN_TITLE)
}

// SetTitle sets the post title.
func (o *postImplementation) SetTitle(title string) PostInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

// GetUpdatedAt returns the last update timestamp as a string.
func (o *postImplementation) GetUpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
// Returns the null datetime if the updated_at field is empty.
func (o *postImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

// SetUpdatedAt sets the last update timestamp.
func (o *postImplementation) SetUpdatedAt(updatedAt string) PostInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

// GetData returns all post data as a map.
func (o *postImplementation) GetData() map[string]string {
	softDeletedAt := o.GetSoftDeletedAt()
	if softDeletedAt == "" {
		softDeletedAt = MAX_DATETIME
	}

	var createdAt, updatedAt string
	if !o.CreatedAtField.CreatedAt.IsZero() {
		createdAt = carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString(carbon.UTC)
	}
	if !o.UpdatedAtField.UpdatedAt.IsZero() {
		updatedAt = carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString(carbon.UTC)
	}

	return map[string]string{
		COLUMN_ID:               o.ShortID.ID,
		COLUMN_AUTHOR_ID:        o.AuthorIDField,
		COLUMN_CANONICAL_URL:    o.CanonicalURLField,
		COLUMN_CONTENT:          o.ContentField,
		COLUMN_FEATURED:         o.FeaturedField,
		COLUMN_IMAGE_URL:        o.ImageURLField,
		COLUMN_MEMO:             o.MemoField,
		COLUMN_META_DESCRIPTION: o.MetaDescriptionField,
		COLUMN_META_KEYWORDS:    o.MetaKeywordsField,
		COLUMN_META_ROBOTS:      o.MetaRobotsField,
		COLUMN_METAS:            o.MetasField,
		COLUMN_PUBLISHED_AT:     carbon.CreateFromStdTime(o.PublishedAtField).ToDateTimeString(carbon.UTC),
		COLUMN_STATUS:           o.StatusField,
		COLUMN_SLUG:             o.SlugField,
		COLUMN_SUMMARY:          o.SummaryField,
		COLUMN_TITLE:            o.TitleField,
		COLUMN_CREATED_AT:       createdAt,
		COLUMN_UPDATED_AT:       updatedAt,
		COLUMN_SOFT_DELETED_AT:  softDeletedAt,
	}
}

// GetDataChanged returns only the fields that have been modified.
// Since neat ORM traits don't track dirty state, return all fields as changed.
func (o *postImplementation) GetDataChanged() map[string]string {
	return o.GetData()
}

// MarkAsNotDirty clears the dirty state of the post.
// No-op since neat ORM traits don't track dirty state.
func (o *postImplementation) MarkAsNotDirty(columns ...string) {
}

// MarkAsDirty marks the post as dirty.
// No-op since neat ORM traits don't track dirty state.
func (o *postImplementation) MarkAsDirty(columns ...string) {
}

// Get retrieves a value by key from the post data.
func (o *postImplementation) Get(key string) string {
	switch key {
	case COLUMN_ID:
		return o.ID
	case COLUMN_AUTHOR_ID:
		return o.AuthorIDField
	case COLUMN_CANONICAL_URL:
		return o.CanonicalURLField
	case COLUMN_CONTENT:
		return o.ContentField
	case COLUMN_FEATURED:
		return o.FeaturedField
	case COLUMN_IMAGE_URL:
		return o.ImageURLField
	case COLUMN_MEMO:
		return o.MemoField
	case COLUMN_META_DESCRIPTION:
		return o.MetaDescriptionField
	case COLUMN_META_KEYWORDS:
		return o.MetaKeywordsField
	case COLUMN_META_ROBOTS:
		return o.MetaRobotsField
	case COLUMN_METAS:
		return o.MetasField
	case COLUMN_PUBLISHED_AT:
		return carbon.CreateFromStdTime(o.PublishedAtField).ToDateTimeString(carbon.UTC)
	case COLUMN_STATUS:
		return o.StatusField
	case COLUMN_SLUG:
		return o.SlugField
	case COLUMN_SUMMARY:
		return o.SummaryField
	case COLUMN_TITLE:
		return o.TitleField
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
	case COLUMN_SOFT_DELETED_AT:
		return o.GetSoftDeletedAt()
	default:
		return ""
	}
}

// Set stores a value by key in the post data.
func (o *postImplementation) Set(key string, value string) {
	switch key {
	case COLUMN_ID:
		o.ShortID.ID = value
	case COLUMN_AUTHOR_ID:
		o.AuthorIDField = value
	case COLUMN_CANONICAL_URL:
		o.CanonicalURLField = value
	case COLUMN_CONTENT:
		o.ContentField = value
	case COLUMN_FEATURED:
		o.FeaturedField = value
	case COLUMN_IMAGE_URL:
		o.ImageURLField = value
	case COLUMN_MEMO:
		o.MemoField = value
	case COLUMN_META_DESCRIPTION:
		o.MetaDescriptionField = value
	case COLUMN_META_KEYWORDS:
		o.MetaKeywordsField = value
	case COLUMN_META_ROBOTS:
		o.MetaRobotsField = value
	case COLUMN_METAS:
		o.MetasField = value
	case COLUMN_PUBLISHED_AT:
		o.PublishedAtField = carbon.Parse(value).StdTime()
	case COLUMN_STATUS:
		o.StatusField = value
	case COLUMN_SLUG:
		o.SlugField = value
	case COLUMN_SUMMARY:
		o.SummaryField = value
	case COLUMN_TITLE:
		o.TitleField = value
	case COLUMN_CREATED_AT:
		if value != "" {
			o.CreatedAtField.CreatedAt = carbon.Parse(value, carbon.UTC).StdTime()
		}
	case COLUMN_UPDATED_AT:
		if value != "" {
			o.UpdatedAtField.UpdatedAt = carbon.Parse(value, carbon.UTC).StdTime()
		}
	case COLUMN_SOFT_DELETED_AT:
		if value == MAX_DATETIME {
			o.SoftDeletedAt = soft_delete.MaxSoftDeletedAt
		} else {
			o.SoftDeletedAt = carbon.Parse(value, carbon.UTC).StdTime()
		}
	}
}

// Hydrate populates the post with data from a map.
func (o *postImplementation) Hydrate(data map[string]string) {
	for key, value := range data {
		o.Set(key, value)
	}
}

// IsDirty returns true if the post has unsaved changes.
// Always returns false since neat ORM traits don't track dirty state.
func (o *postImplementation) IsDirty() bool {
	return false
}

// ============================ TAXONOMY METHODS ============================

// TermIDs retrieves the term IDs for a specific taxonomy from the post metadata.
func (o *postImplementation) TermIDs(taxonomySlug string) []string {
	metas, err := o.GetMetas()
	if err != nil {
		return []string{}
	}

	key := "term_ids_" + taxonomySlug
	jsonStr := lo.ValueOr(metas, key, "")
	if jsonStr == "" {
		return []string{}
	}

	var ids []string
	if err := json.Unmarshal([]byte(jsonStr), &ids); err != nil {
		return []string{}
	}

	return ids
}

// SetTermIDs stores the term IDs for a specific taxonomy in the post metadata.
func (o *postImplementation) SetTermIDs(taxonomySlug string, termIDs []string) PostInterface {
	metas, err := o.GetMetas()
	if err != nil {
		return o
	}

	key := "term_ids_" + taxonomySlug
	if len(termIDs) == 0 {
		delete(metas, key)
	} else {
		jsonBytes, err := json.Marshal(termIDs)
		if err != nil {
			return o
		}
		metas[key] = string(jsonBytes)
	}

	o.SetMetas(metas)
	return o
}

// CategoryIDs retrieves the category IDs associated with this post.
func (o *postImplementation) CategoryIDs() []string {
	return o.TermIDs(TAXONOMY_CATEGORY)
}

// SetCategoryIDs sets the category IDs for this post.
func (o *postImplementation) SetCategoryIDs(ids []string) PostInterface {
	return o.SetTermIDs(TAXONOMY_CATEGORY, ids)
}

// TagIDs retrieves the tag IDs associated with this post.
func (o *postImplementation) TagIDs() []string {
	return o.TermIDs(TAXONOMY_TAG)
}

// SetTagIDs sets the tag IDs for this post.
func (o *postImplementation) SetTagIDs(ids []string) PostInterface {
	return o.SetTermIDs(TAXONOMY_TAG, ids)
}

// ============================ OLD SLUG METHODS ============================

// GetOldSlugs retrieves the array of historical slugs for redirect purposes.
func (o *postImplementation) GetOldSlugs() []string {
	metas, err := o.GetMetas()
	if err != nil {
		return []string{}
	}

	jsonStr := lo.ValueOr(metas, META_KEY_OLD_SLUGS, "")
	if jsonStr == "" {
		return []string{}
	}

	var slugs []string
	if err := json.Unmarshal([]byte(jsonStr), &slugs); err != nil {
		return []string{}
	}

	return slugs
}

// SetOldSlugs sets the array of historical slugs.
func (o *postImplementation) SetOldSlugs(slugs []string) error {
	metas, err := o.GetMetas()
	if err != nil {
		return err
	}

	if len(slugs) == 0 {
		delete(metas, META_KEY_OLD_SLUGS)
	} else {
		jsonBytes, err := json.Marshal(slugs)
		if err != nil {
			return err
		}
		metas[META_KEY_OLD_SLUGS] = string(jsonBytes)
	}

	return o.SetMetas(metas)
}

// AddOldSlug adds a slug to the old slugs history.
func (o *postImplementation) AddOldSlug(slug string) error {
	if slug == "" {
		return nil
	}

	oldSlugs := o.GetOldSlugs()
	// Avoid duplicates
	for _, s := range oldSlugs {
		if s == slug {
			return nil
		}
	}

	oldSlugs = append(oldSlugs, slug)
	return o.SetOldSlugs(oldSlugs)
}

// BlogNoImageUrl returns a default image URL when no featured image is set.
func BlogNoImageUrl() string {
	// return links.NewWebsiteLinks().Cdn("/blogs/default_blog.jpg", map[string]string{})
	//return config.MediaUrl + "/blogs/default_blog.png"
	return "https://picsum.photos/id/20/200/300"
}
