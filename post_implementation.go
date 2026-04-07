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

// var _ dataobject.DataObjectFluentInterface = (*Post)(nil) // verify it extends the data object interface

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

func NewPostFromExistingData(data map[string]string) PostInterface {
	o := &postImplementation{}
	o.Hydrate(data)
	return o
}

type postImplementation struct {
	dataobject.DataObject
}

// ================================== METHODS ==================================
func (o *postImplementation) GetSlug() string {
	return str.Slugify(o.GetTitle(), '-')
}

func (o *postImplementation) GetEditor() string {
	return o.GetMeta("editor")
}

func (o *postImplementation) SetEditor(editor string) PostInterface {
	o.SetMeta("editor", editor)
	return o
}

func (o *postImplementation) GetContentType() string {
	return o.GetMeta("content_type")
}

func (o *postImplementation) SetContentType(contentType string) PostInterface {
	o.SetMeta("content_type", contentType)
	return o
}

func (o *postImplementation) IsDraft() bool {
	return o.GetStatus() == POST_STATUS_DRAFT
}

func (o *postImplementation) IsPublished() bool {
	return o.GetStatus() == POST_STATUS_PUBLISHED
}

func (o *postImplementation) IsContentMarkdown() bool {
	return o.GetContentType() == POST_CONTENT_TYPE_MARKDOWN
}

func (o *postImplementation) IsContentHtml() bool {
	return o.GetContentType() == POST_CONTENT_TYPE_HTML
}

func (o *postImplementation) IsContentPlainText() bool {
	return o.GetContentType() == POST_CONTENT_TYPE_PLAIN_TEXT
}

func (o *postImplementation) IsContentBlocks() bool {
	return o.GetContentType() == POST_CONTENT_TYPE_BLOCKS
}

func (o *postImplementation) IsTrashed() bool {
	return o.GetStatus() == POST_STATUS_TRASH
}

func (o *postImplementation) IsUnpublished() bool {
	return !o.IsPublished()
}

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

func (o *postImplementation) GetAuthorID() string {
	return o.Get(COLUMN_AUTHOR_ID)
}

func (o *postImplementation) SetAuthorID(authorID string) PostInterface {
	o.Set(COLUMN_AUTHOR_ID, authorID)
	return o
}

func (o *postImplementation) GetCanonicalURL() string {
	return o.Get(COLUMN_CANONICAL_URL)
}

func (o *postImplementation) SetCanonicalURL(canonicalURL string) PostInterface {
	o.Set(COLUMN_CANONICAL_URL, canonicalURL)
	return o
}

func (o *postImplementation) GetContent() string {
	return o.Get(COLUMN_CONTENT)
}

func (o *postImplementation) SetContent(content string) PostInterface {
	o.Set(COLUMN_CONTENT, content)
	return o
}

func (o *postImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *postImplementation) SetCreatedAt(createdAt string) PostInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *postImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return carbon.Parse(sb.NULL_DATETIME)
	}
	return carbon.Parse(createdAt)
}

func (o *postImplementation) GetCreatedAtTime() time.Time {
	createdAt := o.GetCreatedAt()
	if createdAt == "" {
		return time.Time{}
	}
	return carbon.Parse(createdAt).StdTime()
}

func (o *postImplementation) GetSoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *postImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	deletedAt := o.GetSoftDeletedAt()
	if deletedAt == "" {
		return carbon.Parse(sb.NULL_DATETIME)
	}
	return carbon.Parse(deletedAt)
}

func (o *postImplementation) SetSoftDeletedAt(deletedAt string) PostInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return o
}

func (o *postImplementation) GetFeatured() string {
	return o.Get(COLUMN_FEATURED)
}

func (o *postImplementation) SetFeatured(featured string) PostInterface {
	o.Set(COLUMN_FEATURED, featured)
	return o
}

func (o *postImplementation) GetID() string {
	return o.Get(COLUMN_ID)
}
func (o *postImplementation) SetID(id string) PostInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *postImplementation) GetImageUrl() string {
	return o.Get(COLUMN_IMAGE_URL)
}

func (o *postImplementation) SetImageUrl(imageURL string) PostInterface {
	o.Set(COLUMN_IMAGE_URL, imageURL)
	return o
}

func (o *postImplementation) GetImageUrlOrDefault() string {
	return lo.Ternary(o.GetImageUrl() == "", BlogNoImageUrl(), o.GetImageUrl())
}

func (o *postImplementation) GetMemo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *postImplementation) SetMemo(memo string) PostInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *postImplementation) GetMeta(key string) string {
	metas, err := o.GetMetas()

	if err != nil {
		return ""
	}

	return lo.ValueOr(metas, key, "")
}

func (o *postImplementation) SetMeta(key string, value string) error {
	metas, err := o.GetMetas()

	if err != nil {
		return err
	}

	metas[key] = value
	return o.SetMetas(metas)
}

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

func (o *postImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	o.Set(COLUMN_METAS, string(mapString))
	return nil
}

func (o *postImplementation) GetMetaDescription() string {
	return o.Get(COLUMN_META_DESCRIPTION)
}

func (o *postImplementation) SetMetaDescription(metaDescription string) PostInterface {
	o.Set(COLUMN_META_DESCRIPTION, metaDescription)
	return o
}

func (o *postImplementation) GetMetaKeywords() string {
	return o.Get(COLUMN_META_KEYWORDS)
}

func (o *postImplementation) SetMetaKeywords(metaKeywords string) PostInterface {
	o.Set(COLUMN_META_KEYWORDS, metaKeywords)
	return o
}

func (o *postImplementation) GetMetaRobots() string {
	return o.Get(COLUMN_META_ROBOTS)
}

func (o *postImplementation) SetMetaRobots(metaRobots string) PostInterface {
	o.Set(COLUMN_META_ROBOTS, metaRobots)
	return o
}

func (o *postImplementation) GetPublishedAt() string {
	return o.Get(COLUMN_PUBLISHED_AT)
}

func (o *postImplementation) SetPublishedAt(status string) PostInterface {
	o.Set(COLUMN_PUBLISHED_AT, status)
	return o
}

func (o *postImplementation) GetPublishedAtCarbon() *carbon.Carbon {
	createdAt := o.GetPublishedAt()
	if createdAt == "" {
		return carbon.Parse(sb.NULL_DATETIME)
	}
	return carbon.Parse(createdAt)
}

func (o *postImplementation) GetPublishedAtTime() time.Time {
	publishedAt := o.GetPublishedAt()
	if publishedAt == "" {
		return time.Time{}
	}
	return carbon.Parse(publishedAt).StdTime()
}

func (o *postImplementation) GetStatus() string {
	return o.Get(COLUMN_STATUS)
}

func (o *postImplementation) SetStatus(status string) PostInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *postImplementation) GetSummary() string {
	return o.Get(COLUMN_SUMMARY)
}

func (o *postImplementation) SetSummary(summary string) PostInterface {
	o.Set(COLUMN_SUMMARY, summary)
	return o
}

func (o *postImplementation) GetTitle() string {
	return o.Get(COLUMN_TITLE)
}

func (o *postImplementation) SetTitle(title string) PostInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

func (o *postImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *postImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := o.GetUpdatedAt()
	if updatedAt == "" {
		return carbon.Parse(sb.NULL_DATETIME)
	}
	return carbon.Parse(updatedAt)
}

func (o *postImplementation) SetUpdatedAt(updatedAt string) PostInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *postImplementation) GetData() map[string]string {
	return o.DataObject.Data()
}

func (o *postImplementation) GetDataChanged() map[string]string {
	return o.DataObject.DataChanged()
}

func (o *postImplementation) MarkAsNotDirty() {
	o.DataObject.MarkAsNotDirty()
}

func (o *postImplementation) Get(key string) string {
	return o.DataObject.Get(key)
}

func (o *postImplementation) Set(key string, value string) {
	o.DataObject.Set(key, value)
}

func (o *postImplementation) Hydrate(data map[string]string) {
	o.DataObject.Hydrate(data)
}

func (o *postImplementation) IsDirty() bool {
	return o.DataObject.IsDirty()
}

func BlogNoImageUrl() string {
	// return links.NewWebsiteLinks().Cdn("/blogs/default_blog.jpg", map[string]string{})
	//return config.MediaUrl + "/blogs/default_blog.png"
	return "https://picsum.photos/id/20/200/300"
}
