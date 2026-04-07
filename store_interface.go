package blogstore

import "context"

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool) StoreInterface
	VersioningEnabled() bool
	TaxonomyEnabled() bool

	PostCount(ctx context.Context, options PostQueryOptions) (int64, error)
	PostCreate(ctx context.Context, post PostInterface) error
	PostDelete(ctx context.Context, post PostInterface) error
	PostDeleteByID(ctx context.Context, postID string) error
	PostFindByID(ctx context.Context, id string) (PostInterface, error)
	PostList(ctx context.Context, options PostQueryOptions) ([]PostInterface, error)
	PostSoftDelete(ctx context.Context, post PostInterface) error
	PostSoftDeleteByID(ctx context.Context, postID string) error
	PostTrash(ctx context.Context, post PostInterface) error
	PostUpdate(ctx context.Context, post PostInterface) error

	// Versioning
	VersioningCreate(ctx context.Context, versioning VersioningInterface) error
	VersioningDelete(ctx context.Context, versioning VersioningInterface) error
	VersioningDeleteByID(ctx context.Context, id string) error
	VersioningFindByID(ctx context.Context, versioningID string) (VersioningInterface, error)
	VersioningList(ctx context.Context, query VersioningQueryInterface) ([]VersioningInterface, error)
	VersioningSoftDelete(ctx context.Context, versioning VersioningInterface) error
	VersioningSoftDeleteByID(ctx context.Context, id string) error
	VersioningUpdate(ctx context.Context, versioning VersioningInterface) error

	// Taxonomy management
	TaxonomyCount(ctx context.Context, options TaxonomyQueryOptions) (int64, error)
	TaxonomyCreate(ctx context.Context, taxonomy TaxonomyInterface) error
	TaxonomyDelete(ctx context.Context, taxonomy TaxonomyInterface) error
	TaxonomyFindByID(ctx context.Context, id string) (TaxonomyInterface, error)
	TaxonomyFindBySlug(ctx context.Context, slug string) (TaxonomyInterface, error)
	TaxonomyList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyInterface, error)
	TaxonomyUpdate(ctx context.Context, taxonomy TaxonomyInterface) error

	// Term management
	TermCount(ctx context.Context, options TermQueryOptions) (int64, error)
	TermCreate(ctx context.Context, term TermInterface) error
	TermDelete(ctx context.Context, term TermInterface) error
	TermFindByID(ctx context.Context, id string) (TermInterface, error)
	TermFindBySlug(ctx context.Context, taxonomySlug, termSlug string) (TermInterface, error)
	TermList(ctx context.Context, options TermQueryOptions) ([]TermInterface, error)
	TermUpdate(ctx context.Context, term TermInterface) error

	// Post-term relationships
	PostTermAdd(ctx context.Context, postID string, termID string, sequence int) error
	PostTermRemove(ctx context.Context, postID string, termID string) error
	PostTerms(ctx context.Context, postID string, taxonomySlug string) ([]TermInterface, error)
	PostSetTerms(ctx context.Context, postID string, taxonomySlug string, termIDs []string) error

	// Utility queries
	TermIncrementCount(ctx context.Context, termID string) error
	TermDecrementCount(ctx context.Context, termID string) error
}
