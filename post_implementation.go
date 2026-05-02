package blogstore

import (
	"encoding/json"
	"time"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/str"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

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
		SetPublishedAt(sb.NULL_DATETIME).
		SetSummary("").
		SetTitle("").
		SetPublishedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME).
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
// It embeds dataobject.DataObject for data storage and change tracking.
type postImplementation struct {
	dataobject.DataObject
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
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the creation timestamp.
func (o *postImplementation) SetCreatedAt(createdAt string) PostInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// GetCreatedAtCarbon returns the creation timestamp as a carbon.Carbon instance.
// Returns the null datetime if the created_at field is empty.
func (o *postImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return carbon.Parse(sb.NULL_DATETIME)
	}
	return carbon.Parse(createdAt)
}

// GetCreatedAtTime returns the creation timestamp as a time.Time instance.
// Returns zero time if the created_at field is empty.
func (o *postImplementation) GetCreatedAtTime() time.Time {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return time.Time{}
	}
	return carbon.Parse(createdAt).StdTime()
}

// GetSoftDeletedAt returns the soft deletion timestamp as a string.
func (o *postImplementation) GetSoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a carbon.Carbon instance.
// Returns the null datetime if the soft_deleted_at field is empty.
func (o *postImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	deletedAt := o.GetSoftDeletedAt()
	if deletedAt == "" {
		return carbon.Parse(sb.NULL_DATETIME)
	}
	return carbon.Parse(deletedAt)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (o *postImplementation) SetSoftDeletedAt(deletedAt string) PostInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
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
		return carbon.Parse(sb.NULL_DATETIME)
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
	return o.Get(COLUMN_UPDATED_AT)
}

// GetUpdatedAtCarbon returns the last update timestamp as a carbon.Carbon instance.
// Returns the null datetime if the updated_at field is empty.
func (o *postImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return carbon.Parse(sb.NULL_DATETIME)
	}
	return carbon.Parse(updatedAt)
}

// SetUpdatedAt sets the last update timestamp.
func (o *postImplementation) SetUpdatedAt(updatedAt string) PostInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

// GetData returns all post data as a map.
func (o *postImplementation) GetData() map[string]string {
	return o.DataObject.Data()
}

// GetDataChanged returns only the fields that have been modified.
func (o *postImplementation) GetDataChanged() map[string]string {
	return o.DataObject.DataChanged()
}

// MarkAsNotDirty clears the dirty state of the post.
func (o *postImplementation) MarkAsNotDirty() {
	o.DataObject.MarkAsNotDirty()
}

// Get retrieves a value by key from the post data.
func (o *postImplementation) Get(key string) string {
	return o.DataObject.Get(key)
}

// Set stores a value by key in the post data.
func (o *postImplementation) Set(key string, value string) {
	o.DataObject.Set(key, value)
}

// Hydrate populates the post with data from a map.
func (o *postImplementation) Hydrate(data map[string]string) {
	o.DataObject.Hydrate(data)
}

// IsDirty returns true if the post has unsaved changes.
func (o *postImplementation) IsDirty() bool {
	return o.DataObject.IsDirty()
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

const META_KEY_OLD_SLUGS = "_wp_old_slug"

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
