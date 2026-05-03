package blogstore

// PostQueryOptions defines the query parameters for retrieving posts.
// These options allow filtering, sorting, and pagination of post results.
type PostQueryOptions struct {
	// ID filters by a single post ID.
	ID string
	// IDIn filters by multiple post IDs.
	IDIn []string
	// Status filters by post status (draft, published, trash, etc.).
	Status string
	// StatusIn filters by multiple post statuses.
	StatusIn []string
	// Slug filters by the post slug.
	Slug string
	// OldSlug filters posts where the old slugs array contains this value.
	OldSlug string
	// Search performs a case-insensitive search on title and content.
	Search string
	// CreatedAtLessThan filters posts created before this timestamp.
	CreatedAtLessThan string
	// CreatedAtGreaterThan filters posts created after this timestamp.
	CreatedAtGreaterThan string
	// Offset is the number of records to skip for pagination.
	Offset int
	// Limit is the maximum number of records to return.
	Limit int
	// SortOrder is the sort direction (asc or desc).
	SortOrder string
	// OrderBy is the field to sort by.
	OrderBy string
	// CountOnly returns only the count, not the actual records.
	CountOnly bool
	// WithDeleted includes soft-deleted posts in the results.
	WithDeleted bool
	// MetaEquals filters posts where the meta JSON column has the specified key-value pair (equality).
	// Example: MetaEquals: map[string]string{"content_type": "plain_text"}
	MetaEquals map[string]string
	// MetaArrayContains filters posts where the meta JSON column's array field contains the specified value.
	// Example: MetaArrayContains: map[string]string{"_old_slugs": "11"}
	MetaArrayContains map[string]string
}
