package blogstore

type PostQueryOptions struct {
	ID                   string
	IDIn                 []string
	Status               string
	StatusIn             []string
	Search               string
	CreatedAtLessThan    string
	CreatedAtGreaterThan string
	Offset               int
	Limit                int
	SortOrder            string
	OrderBy              string
	CountOnly            bool
	WithDeleted          bool
	// MetaContains filters posts where the meta JSON column contains the specified key-value pair.
	// Example: MetaContains: map[string]string{"wp_id": "123"}
	MetaContains map[string]string
}
