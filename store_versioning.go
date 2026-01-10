package blogstore

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dracory/sb"
	"github.com/dracory/versionstore"
)

type versioningMarshalToInterface interface {
	MarshalToVersioning() (string, error)
}

type versioningDataInterface interface {
	Data() map[string]string
}

func (store *store) versioningContentFromEntity(entity any) (string, error) {
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
	for k, v := range d.Data() {
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

func (store *store) versioningCreateIfChanged(ctx context.Context, entityType string, entityID string, content string) error {
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

func (store *store) versioningTrackEntity(ctx context.Context, entityType string, entityID string, entity any) error {
	if !store.VersioningEnabled() {
		return nil
	}

	content, err := store.versioningContentFromEntity(entity)
	if err != nil {
		return err
	}

	return store.versioningCreateIfChanged(ctx, entityType, entityID, content)
}

func (store *store) VersioningCreate(ctx context.Context, version VersioningInterface) error {
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionCreate(store.toQueryableContext(ctx), version)
}

func (store *store) VersioningDelete(ctx context.Context, version VersioningInterface) error {
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionDelete(store.toQueryableContext(ctx), version)
}

func (store *store) VersioningDeleteByID(ctx context.Context, id string) error {
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionDeleteByID(store.toQueryableContext(ctx), id)
}

func (store *store) VersioningFindByID(ctx context.Context, versioningID string) (VersioningInterface, error) {
	if store.versioningStore == nil {
		return nil, nil
	}
	return store.versioningStore.VersionFindByID(store.toQueryableContext(ctx), versioningID)
}

func (store *store) VersioningList(ctx context.Context, query VersioningQueryInterface) ([]VersioningInterface, error) {
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

func (store *store) VersioningSoftDelete(ctx context.Context, versioning VersioningInterface) error {
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionSoftDelete(store.toQueryableContext(ctx), versioning)
}

func (store *store) VersioningSoftDeleteByID(ctx context.Context, id string) error {
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionSoftDeleteByID(store.toQueryableContext(ctx), id)
}

func (store *store) VersioningUpdate(ctx context.Context, version VersioningInterface) error {
	if store.versioningStore == nil {
		return nil
	}
	return store.versioningStore.VersionUpdate(store.toQueryableContext(ctx), version)
}
