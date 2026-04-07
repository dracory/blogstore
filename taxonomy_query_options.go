package blogstore

// TaxonomyQueryOptions defines query options for listing taxonomies
type TaxonomyQueryOptions struct {
	ID        string
	Slug      string
	Search    string
	Limit     int
	Offset    int
	OrderBy   string
	SortOrder string
	CountOnly bool
}

// TermQueryOptions defines query options for listing terms
type TermQueryOptions struct {
	ID           string
	IDIn         []string
	TaxonomyID   string
	TaxonomySlug string
	ParentID     string
	Search       string
	Limit        int
	Offset       int
	OrderBy      string
	SortOrder    string
	CountOnly    bool
}

// TermRelationQueryOptions defines query options for listing term relations
type TermRelationQueryOptions struct {
	PostID       string
	TermID       string
	TaxonomyID   string
	TaxonomySlug string
}
