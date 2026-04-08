package blogstore

import "context"

// StoreInterface defines the complete interface for blog post storage operations,
// including post management, taxonomy/term handling, and optional versioning support.
type StoreInterface interface {
	// AutoMigrate creates or updates the database schema to match the current model definitions.
	// Returns an error if migration fails.
	AutoMigrate() error

	// EnableDebug toggles debug mode logging for database operations.
	// Returns the StoreInterface to allow method chaining.
	EnableDebug(debug bool) StoreInterface

	// VersioningEnabled returns true if versioning support is enabled for this store.
	VersioningEnabled() bool

	// TaxonomyEnabled returns true if taxonomy support is enabled for this store.
	TaxonomyEnabled() bool

	// PostCount returns the total number of posts matching the provided query options.
	// Uses PostQueryOptions to filter by status, type, or other criteria.
	PostCount(ctx context.Context, options PostQueryOptions) (int64, error)

	// PostCreate inserts a new post into the store.
	// Returns an error if the post cannot be created (e.g., duplicate ID or validation failure).
	PostCreate(ctx context.Context, post PostInterface) error

	// PostDelete permanently removes a post from the store.
	// This is a hard delete that cannot be undone.
	PostDelete(ctx context.Context, post PostInterface) error

	// PostDeleteByID permanently removes a post by its ID.
	// Returns an error if the post does not exist.
	PostDeleteByID(ctx context.Context, postID string) error

	// PostFindByID retrieves a post by its unique identifier.
	// Returns the post and nil error on success, or nil and an error if not found.
	PostFindByID(ctx context.Context, id string) (PostInterface, error)

	// PostList retrieves a list of posts matching the provided query options.
	// Supports pagination, sorting, and filtering through PostQueryOptions.
	PostList(ctx context.Context, options PostQueryOptions) ([]PostInterface, error)

	// PostSoftDelete marks a post as deleted without removing it from the database.
	// The post can be restored later. Requires versioning to be enabled.
	PostSoftDelete(ctx context.Context, post PostInterface) error

	// PostSoftDeleteByID marks a post as deleted by ID without removing it from the database.
	// Returns an error if the post does not exist.
	PostSoftDeleteByID(ctx context.Context, postID string) error

	// PostTrash moves a post to the trash status.
	// Trashed posts are not visible in normal queries but can be restored.
	PostTrash(ctx context.Context, post PostInterface) error

	// PostUpdate modifies an existing post in the store.
	// Returns an error if the post does not exist or validation fails.
	PostUpdate(ctx context.Context, post PostInterface) error

	// Versioning methods manage historical versions of posts.

	// VersioningCreate saves a new version record for a post.
	VersioningCreate(ctx context.Context, versioning VersioningInterface) error

	// VersioningDelete permanently removes a versioning record.
	VersioningDelete(ctx context.Context, versioning VersioningInterface) error

	// VersioningDeleteByID permanently removes a versioning record by ID.
	VersioningDeleteByID(ctx context.Context, id string) error

	// VersioningFindByID retrieves a specific version record by ID.
	VersioningFindByID(ctx context.Context, versioningID string) (VersioningInterface, error)

	// VersioningList retrieves version records matching the provided query.
	VersioningList(ctx context.Context, query VersioningQueryInterface) ([]VersioningInterface, error)

	// VersioningSoftDelete marks a version record as deleted without permanent removal.
	VersioningSoftDelete(ctx context.Context, versioning VersioningInterface) error

	// VersioningSoftDeleteByID marks a version record as deleted by ID.
	VersioningSoftDeleteByID(ctx context.Context, id string) error

	// VersioningUpdate modifies an existing version record.
	VersioningUpdate(ctx context.Context, versioning VersioningInterface) error

	// Taxonomy methods manage classification systems (e.g., categories, tags).

	// TaxonomyCount returns the number of taxonomies matching the query options.
	TaxonomyCount(ctx context.Context, options TaxonomyQueryOptions) (int64, error)

	// TaxonomyCreate inserts a new taxonomy into the store.
	TaxonomyCreate(ctx context.Context, taxonomy TaxonomyInterface) error

	// TaxonomyDelete permanently removes a taxonomy from the store.
	TaxonomyDelete(ctx context.Context, taxonomy TaxonomyInterface) error

	// TaxonomyFindByID retrieves a taxonomy by its unique identifier.
	TaxonomyFindByID(ctx context.Context, id string) (TaxonomyInterface, error)

	// TaxonomyFindBySlug retrieves a taxonomy by its URL-friendly slug.
	TaxonomyFindBySlug(ctx context.Context, slug string) (TaxonomyInterface, error)

	// TaxonomyList retrieves taxonomies matching the provided query options.
	TaxonomyList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyInterface, error)

	// TaxonomyUpdate modifies an existing taxonomy.
	TaxonomyUpdate(ctx context.Context, taxonomy TaxonomyInterface) error

	// Term methods manage individual terms within taxonomies.

	// TermCount returns the number of terms matching the query options.
	TermCount(ctx context.Context, options TermQueryOptions) (int64, error)

	// TermCreate inserts a new term into the store.
	TermCreate(ctx context.Context, term TermInterface) error

	// TermDelete permanently removes a term from the store.
	TermDelete(ctx context.Context, term TermInterface) error

	// TermFindByID retrieves a term by its unique identifier.
	TermFindByID(ctx context.Context, id string) (TermInterface, error)

	// TermFindBySlug retrieves a term by its taxonomy slug and term slug.
	TermFindBySlug(ctx context.Context, taxonomySlug, termSlug string) (TermInterface, error)

	// TermList retrieves terms matching the provided query options.
	TermList(ctx context.Context, options TermQueryOptions) ([]TermInterface, error)

	// TermUpdate modifies an existing term.
	TermUpdate(ctx context.Context, term TermInterface) error

	// Post-term relationship methods manage associations between posts and terms.

	// PostAddTerm appends a term to a post (adds at the end of the sequence).
	// Automatically calculates the next available sequence number.
	// Returns an error if taxonomy features are not enabled.
	PostAddTerm(ctx context.Context, postID string, termID string) error

	// PostInsertTermAt associates a term with a post at the specified sequence/order.
	PostInsertTermAt(ctx context.Context, postID string, termID string, sequence int) error

	// PostMoveTermTo moves a term to a specific sequence position on a post.
	// Reorders existing terms by pushing subsequent terms down (incrementing their sequence).
	// Returns an error if the term is not associated with the post.
	PostMoveTermTo(ctx context.Context, postID string, termID string, sequence int) error

	// PostRemoveTerm dissociates a term from a post.
	PostRemoveTerm(ctx context.Context, postID string, termID string) error

	// TermListByPostID retrieves all terms associated with a post for a specific taxonomy.
	TermListByPostID(ctx context.Context, postID string, taxonomySlug string) ([]TermInterface, error)

	// PostSetTerms replaces all terms for a post within a taxonomy with the provided term IDs.
	PostSetTerms(ctx context.Context, postID string, taxonomySlug string, termIDs []string) error

	// PostListByTermID retrieves all posts associated with a specific term ID.
	PostListByTermID(ctx context.Context, termID string, options PostQueryOptions) ([]PostInterface, error)

	// Utility methods provide helper operations for term management.

	// TermIncrementCount increases the usage count for a term.
	TermIncrementCount(ctx context.Context, termID string) error

	// TermDecrementCount decreases the usage count for a term.
	TermDecrementCount(ctx context.Context, termID string) error
}
