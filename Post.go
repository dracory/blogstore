package blogstore

import (
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

// var _ dataobject.DataObjectFluentInterface = (*Post)(nil) // verify it extends the data object interface

func NewPost() *Post {
	o := &Post{}
	o.SetID(uid.HumanUid()).
		SetAuthorID("").
		SetCanonicalURL("").
		SetContent("").
		SetFeatured(NO).
		SetImageUrl("").
		SetMetaDescription("").
		SetMetaKeywords("").
		SetMetaRobots("").
		SetStatus(POST_STATUS_DRAFT).
		SetPublishedAt(sb.NULL_DATETIME).
		SetSummary("").
		SetTitle("").
		SetPublishedAt(carbon.NewCarbon().Now().Format("Y-m-d H:i:s")).
		SetCreatedAt(carbon.NewCarbon().Now().Format("Y-m-d H:i:s")).
		SetUpdatedAt(carbon.NewCarbon().Now().Format("Y-m-d H:i:s")).
		SetDeletedAt(sb.NULL_DATETIME).
		SetMetas(map[string]string{})

	return o
}

func NewPostFromExistingData(data map[string]string) *Post {
	o := &Post{}
	o.Hydrate(data)
	return o
}

type Post struct {
	dataobject.DataObject
}

// ================================== METHODS ==================================
func (o *Post) Slug() string {
	return utils.StrSlugify(o.Title(), '-')
}

// func (o *Post) URL() string {
// 	return links.NewWebsiteLinks().Post(o.ID(), o.Slug())
// }

func (o *Post) Editor() string {
	return o.Meta("editor")
}

// ============================ SETTERS AND GETTERS ============================

func (o *Post) AddMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

func (o *Post) AuthorID() string {
	return o.Get("author_id")
}

func (o *Post) SetAuthorID(authorID string) *Post {
	o.Set("author_id", authorID)
	return o
}

func (o *Post) CanonicalURL() string {
	return o.Get("canonical_url")
}

func (o *Post) CreatedAt() string {
	return o.Get("created_at")
}

func (o *Post) CreatedAtCarbon() carbon.Carbon {
	createdAt := o.CreatedAt()
	if createdAt == "" {
		return carbon.Parse(sb.NULL_DATETIME)
	}
	return carbon.Parse(createdAt)
}

func (o *Post) CreatedAtTime() time.Time {
	createdAt := o.CreatedAt()
	if createdAt == "" {
		return time.Time{}
	}
	return carbon.Parse(createdAt).ToStdTime()
}

func (o *Post) DeletedAt() string {
	return o.Get("deleted_at")
}

func (o *Post) DeletedAtCarbon() carbon.Carbon {
	deletedAt := o.DeletedAt()
	if deletedAt == "" {
		return carbon.Parse(sb.NULL_DATETIME)
	}
	return carbon.Parse(deletedAt)
}

func (o *Post) SetDeletedAt(deletedAt string) *Post {
	o.Set("deleted_at", deletedAt)
	return o
}

func (o *Post) Content() string {
	return o.Get("content")
}

func (o *Post) Featured() string {
	return o.Get("featured")
}

func (o *Post) ID() string {
	return o.Get("id")
}

func (o *Post) IsPublished() bool {
	return o.Status() == POST_STATUS_PUBLISHED
}

func (o *Post) IsTrashed() bool {
	return o.Status() == POST_STATUS_TRASH
}

func (o *Post) IsUnpublished() bool {
	return !o.IsPublished()
}

func (o *Post) ImageUrl() string {
	return o.Get("image_url")
}

func (o *Post) ImageUrlOrDefault() string {
	return lo.Ternary(o.ImageUrl() == "", BlogNoImageUrl(), o.ImageUrl())
}

func (o *Post) MetaDescription() string {
	return o.Get("meta_description")
}

func (o *Post) MetaKeywords() string {
	return o.Get("meta_keywords")
}

func (o *Post) MetaRobots() string {
	return o.Get("meta_robots")
}

func (o *Post) PublishedAt() string {
	return o.Get("published_at")
}

func (o *Post) PublishedAtCarbon() carbon.Carbon {
	createdAt := o.PublishedAt()
	if createdAt == "" {
		return carbon.Parse(sb.NULL_DATETIME)
	}
	return carbon.Parse(createdAt)
}

func (o *Post) PublishedAtTime() time.Time {
	publishedAt := o.PublishedAt()
	if publishedAt == "" {
		return time.Time{}
	}
	return carbon.Parse(publishedAt).ToStdTime()
}

func (o *Post) Status() string {
	return o.Get("status")
}

func (o *Post) Summary() string {
	return o.Get("summary")
}

func (o *Post) Title() string {
	return o.Get("title")
}

func (o *Post) UpdatedAt() string {
	return o.Get("updated_at")
}

func (o *Post) UpdatedAtCarbon() carbon.Carbon {
	updatedAt := o.UpdatedAt()
	if updatedAt == "" {
		return carbon.Parse(sb.NULL_DATETIME)
	}
	return carbon.Parse(updatedAt)
}

func (o *Post) SetUpdatedAt(updatedAt string) *Post {
	o.Set("updated_at", updatedAt)
	return o
}

func (o *Post) SetCanonicalURL(canonicalURL string) *Post {
	o.Set("canonical_url", canonicalURL)
	return o
}

func (o *Post) SetCreatedAt(createdAt string) *Post {
	o.Set("created_at", createdAt)
	return o
}

func (o *Post) SetContent(content string) *Post {
	o.Set("content", content)
	return o
}

func (o *Post) SetFeatured(featured string) *Post {
	o.Set("featured", featured)
	return o
}

func (o *Post) SetID(id string) *Post {
	o.Set("id", id)
	return o
}

func (o *Post) SetImageUrl(imageURL string) *Post {
	o.Set("image_url", imageURL)
	return o
}

func (o *Post) SetMetaDescription(metaDescription string) *Post {
	o.Set("meta_description", metaDescription)
	return o
}

func (o *Post) SetMetaKeywords(metaKeywords string) *Post {
	o.Set("meta_keywords", metaKeywords)
	return o
}

func (o *Post) SetMetaRobots(metaRobots string) *Post {
	o.Set("meta_robots", metaRobots)
	return o
}

func (o *Post) Meta(key string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	return lo.ValueOr(metas, key, "")
}

func (o *Post) Metas() (map[string]string, error) {
	metasStr := o.Get("metas")

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson, errJson := utils.FromJSON(metasStr, map[string]string{})
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return maputils.MapStringAnyToMapStringString(metasJson.(map[string]any)), nil
}

func (o *Post) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}
	o.Set("metas", mapString)
	return nil
}

func (o *Post) SetPublishedAt(status string) *Post {
	o.Set("published_at", status)
	return o
}

func (o *Post) SetStatus(status string) *Post {
	o.Set("status", status)
	return o
}

func (o *Post) SetSummary(summary string) *Post {
	o.Set("summary", summary)
	return o
}

func (o *Post) SetTitle(title string) *Post {
	o.Set("title", title)
	return o
}

func BlogNoImageUrl() string {
	// return links.NewWebsiteLinks().Cdn("/blogs/default_blog.jpg", map[string]string{})
	//return config.MediaUrl + "/blogs/default_blog.png"
	return "https://picsum.photos/id/20/200/300"
}
