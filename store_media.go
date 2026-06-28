package blogstore

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

// MediaCreate inserts a new media into the database.
// Sets created_at and updated_at timestamps automatically.
func (store *storeImplementation) MediaCreate(ctx context.Context, media MediaInterface) error {
	if store.mediaTableName == "" {
		return errors.New("blogstore: media table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if media == nil {
		return errors.New("media is nil")
	}
	if media.GetID() == "" {
		media.SetID(GenerateShortID())
	}
	if media.GetEntityID() == "" {
		return errors.New("media entity_id is empty")
	}

	media.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	media.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	if media.GetSoftDeletedAt() == "" {
		media.SetSoftDeletedAt(MAX_DATETIME)
	}

	metas, _ := media.GetMetas()
	metasJSON := ""
	if len(metas) > 0 {
		metasBytes, err := json.Marshal(metas)
		if err != nil {
			return err
		}
		metasJSON = string(metasBytes)
	}

	row := map[string]any{
		COLUMN_ID:              media.GetID(),
		COLUMN_ENTITY_ID:       media.GetEntityID(),
		COLUMN_TITLE:           media.GetTitle(),
		COLUMN_DESCRIPTION:     media.GetDescription(),
		COLUMN_MEMO:            media.GetMemo(),
		COLUMN_MEDIA_URL:       media.GetURL(),
		COLUMN_MEDIA_TYPE:      media.GetType(),
		COLUMN_FILE_SIZE:       media.GetSize(),
		COLUMN_FILE_EXTENSION:  media.GetExtension(),
		COLUMN_SEQUENCE:        media.GetSequence(),
		COLUMN_STATUS:          media.GetStatus(),
		COLUMN_METAS:           metasJSON,
		COLUMN_CREATED_AT:      media.GetCreatedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:      media.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: media.GetSoftDeletedAtCarbon().StdTime(),
	}

	return store.db.Query().Table(store.mediaTableName).Create(row)
}

// MediaCount returns the total number of media matching the given query options.
func (store *storeImplementation) MediaCount(ctx context.Context, options MediaQueryOptions) (int64, error) {
	if store.mediaTableName == "" {
		return 0, errors.New("blogstore: media table name is empty")
	}
	if ctx == nil {
		return 0, errors.New("ctx is nil")
	}

	q := store.buildMediaQuery(options)

	var count int64
	err := q.Table(store.mediaTableName).Count(&count)
	return count, err
}

// MediaDelete permanently removes a media from the database.
func (store *storeImplementation) MediaDelete(ctx context.Context, media MediaInterface) error {
	if store.mediaTableName == "" {
		return errors.New("blogstore: media table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if media == nil {
		return errors.New("media is nil")
	}

	return store.MediaDeleteByID(ctx, media.GetID())
}

// MediaDeleteByID permanently removes a media by its ID.
func (store *storeImplementation) MediaDeleteByID(ctx context.Context, id string) error {
	if store.mediaTableName == "" {
		return errors.New("blogstore: media table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("media id is empty")
	}

	_, err := store.db.Query().
		Table(store.mediaTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()

	return err
}

// MediaFindByID retrieves a media by its ID.
func (store *storeImplementation) MediaFindByID(ctx context.Context, id string) (MediaInterface, error) {
	if store.mediaTableName == "" {
		return nil, errors.New("blogstore: media table name is empty")
	}
	if id == "" {
		return nil, errors.New("media id is empty")
	}

	list, err := store.MediaList(ctx, MediaQueryOptions{
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

// MediaList retrieves a list of media matching the given query options.
func (store *storeImplementation) MediaList(ctx context.Context, options MediaQueryOptions) ([]MediaInterface, error) {
	if store.mediaTableName == "" {
		return nil, errors.New("blogstore: media table name is empty")
	}
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	type mediaRow struct {
		ID            string    `db:"id"`
		EntityID      string    `db:"entity_id"`
		Title         string    `db:"title"`
		Description   string    `db:"description"`
		Memo          string    `db:"memo"`
		URL           string    `db:"media_url"`
		Type          string    `db:"media_type"`
		Size          string    `db:"file_size"`
		Extension     string    `db:"file_extension"`
		Sequence      int       `db:"sequence"`
		Status        string    `db:"status"`
		Metas         string    `db:"metas"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
		SoftDeletedAt time.Time `db:"soft_deleted_at"`
	}

	q := store.buildMediaQuery(options)

	var rows []mediaRow
	if err := q.Table(store.mediaTableName).Get(&rows); err != nil {
		return []MediaInterface{}, err
	}

	list := make([]MediaInterface, 0, len(rows))
	for _, r := range rows {
		m := &mediaImplementation{
			EntityID:    r.EntityID,
			Title:       r.Title,
			Description: r.Description,
			Memo:        r.Memo,
			URL:         r.URL,
			Type:        r.Type,
			Size:        r.Size,
			Extension:   r.Extension,
			Sequence:    r.Sequence,
			Status:      r.Status,
			Metas:       r.Metas,
		}
		m.ShortID.ID = r.ID
		m.CreatedAt.CreatedAt = r.CreatedAt
		m.UpdatedAt.UpdatedAt = r.UpdatedAt
		m.SoftDeletesMaxDate.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, m)
	}

	return list, nil
}

// MediaListByEntityID retrieves all media attached to a specific entity.
func (store *storeImplementation) MediaListByEntityID(ctx context.Context, entityID string) ([]MediaInterface, error) {
	if store.mediaTableName == "" {
		return nil, errors.New("blogstore: media table name is empty")
	}
	if entityID == "" {
		return nil, errors.New("entity id is empty")
	}

	return store.MediaList(ctx, MediaQueryOptions{
		EntityID:  entityID,
		OrderBy:   COLUMN_SEQUENCE,
		SortOrder: "ASC",
	})
}

// MediaSoftDelete marks a media as deleted by setting the soft_deleted_at timestamp.
func (store *storeImplementation) MediaSoftDelete(ctx context.Context, media MediaInterface) error {
	if store.mediaTableName == "" {
		return errors.New("blogstore: media table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if media == nil {
		return errors.New("media is nil")
	}

	media.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.MediaUpdate(ctx, media)
}

// MediaSoftDeleteByID marks a media as deleted by its ID.
func (store *storeImplementation) MediaSoftDeleteByID(ctx context.Context, id string) error {
	if store.mediaTableName == "" {
		return errors.New("blogstore: media table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("media id is empty")
	}

	media, err := store.MediaFindByID(ctx, id)
	if err != nil {
		return err
	}
	if media == nil {
		return errors.New("media not found")
	}

	return store.MediaSoftDelete(ctx, media)
}

// MediaUpdate updates an existing media in the database.
func (store *storeImplementation) MediaUpdate(ctx context.Context, media MediaInterface) error {
	if store.mediaTableName == "" {
		return errors.New("blogstore: media table name is empty")
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if media == nil {
		return errors.New("media is nil")
	}

	media.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	metas, _ := media.GetMetas()
	metasJSON := ""
	if len(metas) > 0 {
		metasBytes, err := json.Marshal(metas)
		if err != nil {
			return err
		}
		metasJSON = string(metasBytes)
	}

	row := map[string]any{
		COLUMN_ENTITY_ID:       media.GetEntityID(),
		COLUMN_TITLE:           media.GetTitle(),
		COLUMN_DESCRIPTION:     media.GetDescription(),
		COLUMN_MEMO:            media.GetMemo(),
		COLUMN_MEDIA_URL:       media.GetURL(),
		COLUMN_MEDIA_TYPE:      media.GetType(),
		COLUMN_FILE_SIZE:       media.GetSize(),
		COLUMN_FILE_EXTENSION:  media.GetExtension(),
		COLUMN_SEQUENCE:        media.GetSequence(),
		COLUMN_STATUS:          media.GetStatus(),
		COLUMN_METAS:           metasJSON,
		COLUMN_UPDATED_AT:      media.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: media.GetSoftDeletedAtCarbon().StdTime(),
	}

	_, err := store.db.Query().
		Table(store.mediaTableName).
		Where(COLUMN_ID+" = ?", media.GetID()).
		Update(row)

	return err
}

// buildMediaQuery builds a neat query from the media query options.
func (store *storeImplementation) buildMediaQuery(options MediaQueryOptions) contractsorm.Query {
	q := store.db.Query().Table(store.mediaTableName)

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

	if options.EntityID != "" {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID)
	}

	if len(options.EntityIDIn) > 0 {
		inClause := COLUMN_ENTITY_ID + " IN ("
		placeholders := make([]interface{}, 0, len(options.EntityIDIn))
		for i, id := range options.EntityIDIn {
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
		q = q.Where(COLUMN_MEDIA_TYPE+" = ?", options.Type)
	}

	if options.Status != "" {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status)
	}

	if options.Search != "" {
		q = q.Where(COLUMN_TITLE+" LIKE ?", "%"+options.Search+"%")
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
