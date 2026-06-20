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

// versioningMarshalToInterface is an internal interface for entities that can serialize to versioning content.
type versioningMarshalToInterface interface {
	MarshalToVersioning() (string, error)
}

// versioningDataInterface is an internal interface for entities that expose their data as a map.
type versioningDataInterface interface {
	GetData() map[string]string
}

// versioningContentFromEntity extracts versionable content from an entity.
// Supports entities implementing versioningMarshalToInterface or versioningDataInterface.
func (store *storeImplementation) versioningContentFromEntity(entity any) (string, error) {
	if entity == nil {
		return "", errors.New("entity is nil")
	}

	if v, ok := entity.(versioningMarshalToInterface); ok {
		return v.MarshalToVersioning()
	}

	d, ok := entity.(versioningDataInterface)
	if !ok {
		return "", errors.New("entity does not support versioning")
	}

	versionedData := map[string]string{}
	for k, v := range d.GetData() {
		if k == COLUMN_CREATED_AT ||
			k == COLUMN_UPDATED_AT ||
			k == COLUMN_SOFT_DELETED_AT {
			continue
		}
		versionedData[k] = v
	}

	b, err := json.Marshal(versionedData)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// versioningCreateIfChanged creates a new version entry only if the content has changed.
// Compares with the most recent version to avoid duplicate entries.
func (store *storeImplementation) versioningCreateIfChanged(ctx context.Context, entityType string, entityID string, content string) error {
	if !store.VersioningEnabled() {
		return nil
	}

	if store.versioningTableName == "" {
		return errors.New("blogstore: versioning table name is empty")
	}

	if entityType == "" {
		return errors.New("blogstore: entityType is empty")
	}

	if entityID == "" {
		return errors.New("blogstore: entityID is empty")
	}

	lastVersioningList, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(entityType).
		SetEntityID(entityID).
		SetOrderBy(COLUMN_CREATED_AT).
		SetSortOrder("DESC").
		SetLimit(1))
	if err != nil {
		return err
	}

	if len(lastVersioningList) > 0 {
		lastVersioning := lastVersioningList[0]
		if lastVersioning != nil && lastVersioning.Content() == content {
			return nil
		}
	}

	return store.VersioningCreate(ctx, NewVersioning().
		SetEntityID(entityID).
		SetEntityType(entityType).
		SetContent(content))
}

// versioningTrackEntity tracks an entity by creating a version entry if changed.
// This is the main entry point for automatic versioning of posts and other entities.
func (store *storeImplementation) versioningTrackEntity(ctx context.Context, entityType string, entityID string, entity any) error {
	if !store.VersioningEnabled() {
		return nil
	}

	content, err := store.versioningContentFromEntity(entity)
	if err != nil {
		return err
	}

	return store.versioningCreateIfChanged(ctx, entityType, entityID, content)
}

// VersioningCreate creates a new version entry in the versioning store.
func (store *storeImplementation) VersioningCreate(ctx context.Context, version VersioningInterface) error {
	if store.versioningTableName == "" {
		return nil
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if version == nil {
		return errors.New("versioning is nil")
	}
	if version.ID() == "" {
		return errors.New("versioning id is empty")
	}
	if version.EntityType() == "" {
		return errors.New("versioning entity type is empty")
	}
	if version.EntityID() == "" {
		return errors.New("versioning entity id is empty")
	}
	if version.GetCreatedAt() == "" {
		version.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if version.GetSoftDeletedAt() == "" {
		version.SetSoftDeletedAt(MAX_DATETIME)
	}

	row := map[string]any{
		COLUMN_ID:              version.ID(),
		COLUMN_ENTITY_TYPE:     version.EntityType(),
		COLUMN_ENTITY_ID:       version.EntityID(),
		COLUMN_CONTENT:         version.Content(),
		COLUMN_CREATED_AT:      version.GetCreatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: version.GetSoftDeletedAtCarbon().StdTime(),
	}

	return store.db.Query().Table(store.versioningTableName).Create(row)
}

// VersioningDelete permanently removes a version entry from the versioning store.
func (store *storeImplementation) VersioningDelete(ctx context.Context, version VersioningInterface) error {
	if store.versioningTableName == "" {
		return nil
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if version == nil {
		return errors.New("versioning is nil")
	}

	return store.VersioningDeleteByID(ctx, version.ID())
}

// VersioningDeleteByID permanently removes a version entry by its ID.
func (store *storeImplementation) VersioningDeleteByID(ctx context.Context, id string) error {
	if store.versioningTableName == "" {
		return nil
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("versioning id is empty")
	}

	_, err := store.db.Query().
		Table(store.versioningTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

// VersioningFindByID retrieves a version entry by its ID.
func (store *storeImplementation) VersioningFindByID(ctx context.Context, versioningID string) (VersioningInterface, error) {
	if store.versioningTableName == "" {
		return nil, nil
	}
	if versioningID == "" {
		return nil, errors.New("versioning id is empty")
	}

	list, err := store.VersioningList(ctx, NewVersioningQuery().SetID(versioningID).SetLimit(1))
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// VersioningList retrieves a list of version entries matching the given query.
func (store *storeImplementation) VersioningList(ctx context.Context, query VersioningQueryInterface) ([]VersioningInterface, error) {
	if store.versioningTableName == "" {
		return []VersioningInterface{}, nil
	}
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	type versioningRow struct {
		ID            string    `db:"id"`
		EntityType    string    `db:"entity_type"`
		EntityID      string    `db:"entity_id"`
		Content       string    `db:"content"`
		CreatedAt     time.Time `db:"created_at"`
		SoftDeletedAt time.Time `db:"soft_deleted_at"`
	}

	q := store.buildVersioningQuery(query)
	q = q.Table(store.versioningTableName)

	if len(query.Columns()) > 0 {
		q = q.Select(query.Columns())
	}

	var rows []versioningRow
	if err := q.Get(&rows); err != nil {
		return []VersioningInterface{}, err
	}

	list := make([]VersioningInterface, 0, len(rows))
	for _, r := range rows {
		v := &versioningImplementation{
			EntityTypeField: r.EntityType,
			EntityIDField:   r.EntityID,
			ContentField:    r.Content,
			CreatedAt:       r.CreatedAt,
		}
		v.ShortID.ID = r.ID
		v.SoftDeletesMaxDate.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, v)
	}

	return list, nil
}

// VersioningSoftDelete marks a version entry as deleted.
func (store *storeImplementation) VersioningSoftDelete(ctx context.Context, versioning VersioningInterface) error {
	if store.versioningTableName == "" {
		return nil
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if versioning == nil {
		return errors.New("versioning is nil")
	}

	versioning.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.VersioningUpdate(ctx, versioning)
}

// VersioningSoftDeleteByID marks a version entry as deleted by its ID.
func (store *storeImplementation) VersioningSoftDeleteByID(ctx context.Context, id string) error {
	if store.versioningTableName == "" {
		return nil
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("versioning id is empty")
	}

	version, err := store.VersioningFindByID(ctx, id)
	if err != nil {
		return err
	}
	if version == nil {
		return errors.New("versioning not found")
	}

	return store.VersioningSoftDelete(ctx, version)
}

// VersioningUpdate updates an existing version entry in the versioning store.
func (store *storeImplementation) VersioningUpdate(ctx context.Context, version VersioningInterface) error {
	if store.versioningTableName == "" {
		return nil
	}
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if version == nil {
		return errors.New("versioning is nil")
	}

	row := map[string]any{
		COLUMN_SOFT_DELETED_AT: version.GetSoftDeletedAtCarbon().StdTime(),
	}

	_, err := store.db.Query().Table(store.versioningTableName).Where(COLUMN_ID+" = ?", version.ID()).Update(row)
	return err
}

// buildVersioningQuery builds a neat query from the versioning query interface.
func (store *storeImplementation) buildVersioningQuery(options VersioningQueryInterface) contractsorm.Query {
	// Use Model() to enable neat's automatic soft delete handling via SoftDeletesMaxDate
	// Then override the table name since versioningImplementation doesn't implement TableName()
	// Use Select("*") because versioningImplementation wraps timestamps in named struct fields
	// (CreatedAt) which neat's column extractor skips
	q := store.db.Query().Model(&versioningImplementation{}).Select("*")

	if options == nil {
		return q
	}

	if options.HasID() && options.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.HasEntityType() && options.EntityType() != "" {
		q = q.Where(COLUMN_ENTITY_TYPE+" = ?", options.EntityType())
	}

	if options.HasEntityID() && options.EntityID() != "" {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID())
	}

	if options.HasLimit() && options.Limit() > 0 {
		q = q.Limit(options.Limit())
	}

	if options.HasOffset() && options.Offset() > 0 {
		q = q.Offset(int(options.Offset()))
	}

	if options.HasOrderBy() && options.OrderBy() != "" {
		if options.HasSortOrder() && strings.ToLower(options.SortOrder()) == "asc" {
			q = q.OrderBy(options.OrderBy())
		} else {
			q = q.OrderByDesc(options.OrderBy())
		}
	}

	// Active records have soft_deleted_at > NOW (soft-deleted have soft_deleted_at <= NOW)
	if options.HasSoftDeletedIncluded() && options.SoftDeletedIncluded() {
		q = q.WithSoftDeleted()
	} else {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).StdTime())
	}

	return q
}
