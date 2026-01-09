package blogstore

import "context"

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
	for i, v := range list {
		newList[i] = v
	}
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
