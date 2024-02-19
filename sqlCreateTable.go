package blogstore

import (
	"github.com/gouniverse/sb"
)

// SQLCreateTable returns a SQL string for creating the cache table
func (st *Store) sqlCreateTable() string {
	sql := sb.NewBuilder(sb.DatabaseDriverName(st.db)).
		Table(st.postTableName).
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   "status",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   "title",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: "content",
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: "summary",
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: "image_url",
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name:   "featured",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 3,
		}).
		Column(sb.Column{
			Name:   "author_id",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   "canonical_url",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name:   "meta_keywords",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name:   "meta_description",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name:   "meta_robots",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: "metas",
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: "published_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: "updated_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: "deleted_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql
}
