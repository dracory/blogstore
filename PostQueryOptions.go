package blogstore

type PostQueryOptions struct {
	ID                   string
	IDIn                 []string
	Status               string
	StatusIn             []string
	CreatedAtLessThan    string
	CreatedAtGreaterThan string
	Offset               int
	Limit                int
	SortOrder            string
	OrderBy              string
	CountOnly            bool
	WithDeleted          bool
}
