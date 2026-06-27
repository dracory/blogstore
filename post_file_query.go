package blogstore

// PostFileQueryOptions defines query options for listing post files.
type PostFileQueryOptions struct {
	// ID filters by a single post file ID.
	ID string
	// IDIn filters by multiple post file IDs.
	IDIn []string
	// PostID filters by the associated post ID.
	PostID string
	// PostIDIn filters by multiple post IDs.
	PostIDIn []string
	// Extension filters by file extension.
	Extension string
	// Type filters by file type (mime type).
	Type string
	// Search performs a case-insensitive search on name.
	Search string
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
	// WithDeleted includes soft-deleted post files in the results.
	WithDeleted bool
}
