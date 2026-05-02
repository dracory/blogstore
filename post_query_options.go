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
	// MetaContains filters posts where the meta JSON column contains the specified key-value pair.
	// Example: MetaContains: map[string]string{"wp_id": "123"}
	MetaContains map[string]string
}
