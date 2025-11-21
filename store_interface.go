package blogstore

import "context"

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool) StoreInterface

	PostCount(ctx context.Context, options PostQueryOptions) (int64, error)
	PostCreate(ctx context.Context, post *Post) error
	PostDelete(ctx context.Context, post *Post) error
	PostDeleteByID(ctx context.Context, postID string) error
	PostFindByID(ctx context.Context, id string) (*Post, error)
	PostList(ctx context.Context, options PostQueryOptions) ([]Post, error)
	PostSoftDelete(ctx context.Context, post *Post) error
	PostSoftDeleteByID(ctx context.Context, postID string) error
	PostTrash(ctx context.Context, post *Post) error
	PostUpdate(ctx context.Context, post *Post) error
}
