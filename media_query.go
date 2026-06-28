package blogstore

// MediaQueryOptions defines query options for listing media.
type MediaQueryOptions struct {
	// ID filters by a single media ID.
	ID string
	// IDIn filters by multiple media IDs.
	IDIn []string
	// EntityID filters by the associated entity ID.
	EntityID string
	// EntityIDIn filters by multiple entity IDs.
	EntityIDIn []string
	// Extension filters by file extension.
	Extension string
	// Type filters by file type (mime type).
	Type string
	// Status filters by media status.
	Status string
	// Search performs a case-insensitive search on title.
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
	// WithDeleted includes soft-deleted media in the results.
	WithDeleted bool
}
