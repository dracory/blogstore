package blogstore

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dracory/sb"
	"github.com/dracory/versionstore"
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

	if store.versioningStore == nil {
		return errors.New("blogstore: versioning store is nil")
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
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder(sb.DESC).
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
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionCreate(store.toQueryableContext(ctx), version)
}

// VersioningDelete permanently removes a version entry from the versioning store.
func (store *storeImplementation) VersioningDelete(ctx context.Context, version VersioningInterface) error {
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionDelete(store.toQueryableContext(ctx), version)
}

// VersioningDeleteByID permanently removes a version entry by its ID.
func (store *storeImplementation) VersioningDeleteByID(ctx context.Context, id string) error {
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionDeleteByID(store.toQueryableContext(ctx), id)
}

// VersioningFindByID retrieves a version entry by its ID.
func (store *storeImplementation) VersioningFindByID(ctx context.Context, versioningID string) (VersioningInterface, error) {
	if store.versioningStore == nil {
		return nil, nil
	}
	return store.versioningStore.VersionFindByID(store.toQueryableContext(ctx), versioningID)
}

// VersioningList retrieves a list of version entries matching the given query.
func (store *storeImplementation) VersioningList(ctx context.Context, query VersioningQueryInterface) ([]VersioningInterface, error) {
	if store.versioningStore == nil {
		return []VersioningInterface{}, nil
	}
	list, err := store.versioningStore.VersionList(store.toQueryableContext(ctx), query)
	if err != nil {
		return nil, err
	}

	newList := make([]VersioningInterface, len(list))
	copy(newList, list)
	return newList, nil
}

// VersioningSoftDelete marks a version entry as deleted.
func (store *storeImplementation) VersioningSoftDelete(ctx context.Context, versioning VersioningInterface) error {
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionSoftDelete(store.toQueryableContext(ctx), versioning)
}

// VersioningSoftDeleteByID marks a version entry as deleted by its ID.
func (store *storeImplementation) VersioningSoftDeleteByID(ctx context.Context, id string) error {
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionSoftDeleteByID(store.toQueryableContext(ctx), id)
}

// VersioningUpdate updates an existing version entry in the versioning store.
func (store *storeImplementation) VersioningUpdate(ctx context.Context, version VersioningInterface) error {
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionUpdate(store.toQueryableContext(ctx), version)
}
