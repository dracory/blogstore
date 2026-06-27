package blogstore

import (
	"context"
	"errors"
	"strings"
	"time"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

// PostFileCreate inserts a new post file into the database.
// Sets created_at and updated_at timestamps automatically.
func (store *storeImplementation) PostFileCreate(ctx context.Context, file PostFileInterface) error {
	if store.postFileTableName == "" {
		return errors.New("blogstore: post file table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if file == nil {
		return errors.New("post file is nil")
	}
	if file.GetID() == "" {
		file.SetID(GenerateShortID())
	}
	if file.GetPostID() == "" {
		return errors.New("post file post_id is empty")
	}

	file.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	file.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	if file.GetSoftDeletedAt() == "" {
		file.SetSoftDeletedAt(MAX_DATETIME)
	}

	row := map[string]any{
		COLUMN_ID:              file.GetID(),
		COLUMN_POST_ID:         file.GetPostID(),
		COLUMN_NAME:            file.GetName(),
		COLUMN_URL:             file.GetURL(),
		COLUMN_FILE_TYPE:       file.GetType(),
		COLUMN_FILE_SIZE:       file.GetSize(),
		COLUMN_FILE_EXTENSION:  file.GetExtension(),
		COLUMN_SEQUENCE:        file.GetSequence(),
		COLUMN_CREATED_AT:      file.GetCreatedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:      file.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: file.GetSoftDeletedAtCarbon().StdTime(),
	}

	return store.db.Query().Table(store.postFileTableName).Create(row)
}

// PostFileCount returns the total number of post files matching the given query options.
func (store *storeImplementation) PostFileCount(ctx context.Context, options PostFileQueryOptions) (int64, error) {
	if store.postFileTableName == "" {
		return 0, errors.New("blogstore: post file table name is empty")
	}
	if ctx == nil {
		return 0, errors.New("ctx is nil")
	}

	q := store.buildPostFileQuery(options)

	var count int64
	err := q.Table(store.postFileTableName).Count(&count)
	return count, err
}

// PostFileDelete permanently removes a post file from the database.
func (store *storeImplementation) PostFileDelete(ctx context.Context, file PostFileInterface) error {
	if store.postFileTableName == "" {
		return errors.New("blogstore: post file table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if file == nil {
		return errors.New("post file is nil")
	}

	return store.PostFileDeleteByID(ctx, file.GetID())
}

// PostFileDeleteByID permanently removes a post file by its ID.
func (store *storeImplementation) PostFileDeleteByID(ctx context.Context, id string) error {
	if store.postFileTableName == "" {
		return errors.New("blogstore: post file table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("post file id is empty")
	}

	_, err := store.db.Query().
		Table(store.postFileTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()

	return err
}

// PostFileFindByID retrieves a post file by its ID.
func (store *storeImplementation) PostFileFindByID(ctx context.Context, id string) (PostFileInterface, error) {
	if store.postFileTableName == "" {
		return nil, errors.New("blogstore: post file table name is empty")
	}
	if id == "" {
		return nil, errors.New("post file id is empty")
	}

	list, err := store.PostFileList(ctx, PostFileQueryOptions{
		ID:    id,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// PostFileList retrieves a list of post files matching the given query options.
func (store *storeImplementation) PostFileList(ctx context.Context, options PostFileQueryOptions) ([]PostFileInterface, error) {
	if store.postFileTableName == "" {
		return nil, errors.New("blogstore: post file table name is empty")
	}
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	type postFileRow struct {
		ID            string    `db:"id"`
		PostID        string    `db:"post_id"`
		Name          string    `db:"name"`
		URL           string    `db:"url"`
		Type          string    `db:"file_type"`
		Size          string    `db:"file_size"`
		Extension     string    `db:"file_extension"`
		Sequence      int       `db:"sequence"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
		SoftDeletedAt time.Time `db:"soft_deleted_at"`
	}

	q := store.buildPostFileQuery(options)

	var rows []postFileRow
	if err := q.Table(store.postFileTableName).Get(&rows); err != nil {
		return []PostFileInterface{}, err
	}

	list := make([]PostFileInterface, 0, len(rows))
	for _, r := range rows {
		f := &postFileImplementation{
			PostIDField:    r.PostID,
			NameField:      r.Name,
			URLField:       r.URL,
			TypeField:      r.Type,
			SizeField:      r.Size,
			ExtensionField: r.Extension,
			SequenceField:  r.Sequence,
		}
		f.ShortID.ID = r.ID
		f.CreatedAtField.CreatedAt = r.CreatedAt
		f.UpdatedAtField.UpdatedAt = r.UpdatedAt
		f.SoftDeletesMaxDate.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, f)
	}

	return list, nil
}

// PostFileListByPostID retrieves all files attached to a specific post.
func (store *storeImplementation) PostFileListByPostID(ctx context.Context, postID string) ([]PostFileInterface, error) {
	if store.postFileTableName == "" {
		return nil, errors.New("blogstore: post file table name is empty")
	}
	if postID == "" {
		return nil, errors.New("post id is empty")
	}

	return store.PostFileList(ctx, PostFileQueryOptions{
		PostID:    postID,
		OrderBy:   COLUMN_SEQUENCE,
		SortOrder: "ASC",
	})
}

// PostFileSoftDelete marks a post file as deleted by setting the soft_deleted_at timestamp.
func (store *storeImplementation) PostFileSoftDelete(ctx context.Context, file PostFileInterface) error {
	if store.postFileTableName == "" {
		return errors.New("blogstore: post file table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if file == nil {
		return errors.New("post file is nil")
	}

	file.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.PostFileUpdate(ctx, file)
}

// PostFileSoftDeleteByID marks a post file as deleted by its ID.
func (store *storeImplementation) PostFileSoftDeleteByID(ctx context.Context, id string) error {
	if store.postFileTableName == "" {
		return errors.New("blogstore: post file table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("post file id is empty")
	}

	file, err := store.PostFileFindByID(ctx, id)
	if err != nil {
		return err
	}
	if file == nil {
		return errors.New("post file not found")
	}

	return store.PostFileSoftDelete(ctx, file)
}

// PostFileUpdate updates an existing post file in the database.
func (store *storeImplementation) PostFileUpdate(ctx context.Context, file PostFileInterface) error {
	if store.postFileTableName == "" {
		return errors.New("blogstore: post file table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if file == nil {
		return errors.New("post file is nil")
	}

	file.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row := map[string]any{
		COLUMN_POST_ID:         file.GetPostID(),
		COLUMN_NAME:            file.GetName(),
		COLUMN_URL:             file.GetURL(),
		COLUMN_FILE_TYPE:       file.GetType(),
		COLUMN_FILE_SIZE:       file.GetSize(),
		COLUMN_FILE_EXTENSION:  file.GetExtension(),
		COLUMN_SEQUENCE:        file.GetSequence(),
		COLUMN_UPDATED_AT:      file.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: file.GetSoftDeletedAtCarbon().StdTime(),
	}

	_, err := store.db.Query().
		Table(store.postFileTableName).
		Where(COLUMN_ID+" = ?", file.GetID()).
		Update(row)

	return err
}

// buildPostFileQuery builds a neat query from the post file query options.
func (store *storeImplementation) buildPostFileQuery(options PostFileQueryOptions) contractsorm.Query {
	q := store.db.Query().Table(store.postFileTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if len(options.IDIn) > 0 {
		inClause := COLUMN_ID + " IN ("
		placeholders := make([]interface{}, 0, len(options.IDIn))
		for i, id := range options.IDIn {
			if i > 0 {
				inClause += ", "
			}
			inClause += "?"
			placeholders = append(placeholders, id)
		}
		inClause += ")"
		q = q.Where(inClause, placeholders...)
	}

	if options.PostID != "" {
		q = q.Where(COLUMN_POST_ID+" = ?", options.PostID)
	}

	if len(options.PostIDIn) > 0 {
		inClause := COLUMN_POST_ID + " IN ("
		placeholders := make([]interface{}, 0, len(options.PostIDIn))
		for i, id := range options.PostIDIn {
			if i > 0 {
				inClause += ", "
			}
			inClause += "?"
			placeholders = append(placeholders, id)
		}
		inClause += ")"
		q = q.Where(inClause, placeholders...)
	}

	if options.Extension != "" {
		q = q.Where(COLUMN_FILE_EXTENSION+" = ?", options.Extension)
	}

	if options.Type != "" {
		q = q.Where(COLUMN_FILE_TYPE+" = ?", options.Type)
	}

	if options.Search != "" {
		q = q.Where(COLUMN_NAME+" LIKE ?", "%"+options.Search+"%")
	}

	if options.OrderBy != "" {
		order := options.SortOrder
		if order == "" {
			order = "DESC"
		}
		if strings.ToLower(order) == "asc" {
			q = q.OrderBy(options.OrderBy)
		} else {
			q = q.OrderByDesc(options.OrderBy)
		}
	}

	if options.Limit > 0 {
		q = q.Limit(options.Limit)
	}

	if options.Offset > 0 {
		q = q.Offset(options.Offset)
	}

	// Handle soft delete filtering
	if options.WithDeleted {
		q = q.WithSoftDeleted()
	} else {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).StdTime())
	}

	return q
}
