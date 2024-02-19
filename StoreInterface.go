package blogstore

type StoreInterface interface {
	PostCreate(post *Post) error
	PostDelete(post *Post) error
	PostDeleteByID(postID string) error
	PostFindByID(id string) (*Post, error)
	PostList(options PostQueryOptions) ([]Post, error)
	PostUpdate(post *Post) error
}
