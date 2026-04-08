package blogstore

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dracory/sb"
	"github.com/dracory/versionstore"
	_ "modernc.org/sqlite"
)

func TestVersioningContentFromEntity_NilEntity(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	content, err := s.versioningContentFromEntity(nil)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if content != "" {
		t.Errorf("expected empty content, got %q", content)
	}
	if !strings.Contains(err.Error(), "entity is nil") {
		t.Errorf("expected error to contain 'entity is nil', got %q", err.Error())
	}
}

func TestVersioningContentFromEntity_WithMarshalToVersioning(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	post := NewPost()
	post.SetTitle("Test Title").SetContent("Test Content")

	content, err := s.versioningContentFromEntity(post)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if content == "" {
		t.Error("expected non-empty content")
	}

	var data map[string]string
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if data["title"] != "Test Title" {
		t.Errorf("expected title 'Test Title', got %q", data["title"])
	}
	if data["content"] != "Test Content" {
		t.Errorf("expected content 'Test Content', got %q", data["content"])
	}
}

func TestVersioningContentFromEntity_UnsupportedEntity(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	type unsupportedEntity struct{}

	content, err := s.versioningContentFromEntity(&unsupportedEntity{})
	if err == nil {
		t.Error("expected error, got nil")
	}
	if content != "" {
		t.Errorf("expected empty content, got %q", content)
	}
	if !strings.Contains(err.Error(), "entity does not support versioning") {
		t.Errorf("expected error to contain 'entity does not support versioning', got %q", err.Error())
	}
}

func TestVersioningCreateIfChanged_VersioningDisabled(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   false,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()
	err = s.versioningCreateIfChanged(ctx, VERSIONING_TYPE_POST, "post-123", "content")
	if err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestVersioningCreateIfChanged_NilVersioningStore(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	s.versioningStore = nil

	ctx := context.Background()
	err = s.versioningCreateIfChanged(ctx, VERSIONING_TYPE_POST, "post-123", "content")
	if err == nil {
		t.Error("expected error, got nil")
	} else if !strings.Contains(err.Error(), "blogstore: versioning store is nil") {
		t.Errorf("expected error to contain 'blogstore: versioning store is nil', got %q", err.Error())
	}
}

func TestVersioningCreateIfChanged_EmptyEntityType(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()
	err = s.versioningCreateIfChanged(ctx, "", "post-123", "content")
	if err == nil {
		t.Error("expected error, got nil")
	} else if !strings.Contains(err.Error(), "blogstore: entityType is empty") {
		t.Errorf("expected error to contain 'blogstore: entityType is empty', got %q", err.Error())
	}
}

func TestVersioningCreateIfChanged_EmptyEntityID(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()
	err = s.versioningCreateIfChanged(ctx, VERSIONING_TYPE_POST, "", "content")
	if err == nil {
		t.Error("expected error, got nil")
	} else if !strings.Contains(err.Error(), "blogstore: entityID is empty") {
		t.Errorf("expected error to contain 'blogstore: entityID is empty', got %q", err.Error())
	}
}

func TestVersioningCreateIfChanged_NoChange(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()
	entityID := "post-123"
	content := `{"title":"Test Title"}`

	err = s.VersioningCreate(ctx, NewVersioning().
		SetEntityID(entityID).
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(content))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = s.versioningCreateIfChanged(ctx, VERSIONING_TYPE_POST, entityID, content)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	list, err := s.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_POST).
		SetEntityID(entityID))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 versioning record, got %d", len(list))
	}
}

func TestVersioningCreateIfChanged_WithChange(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()
	entityID := "post-456"
	content1 := `{"title":"Test Title 1"}`
	content2 := `{"title":"Test Title 2"}`

	err = s.VersioningCreate(ctx, NewVersioning().
		SetEntityID(entityID).
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(content1))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = s.versioningCreateIfChanged(ctx, VERSIONING_TYPE_POST, entityID, content2)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	list, err := s.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_POST).
		SetEntityID(entityID).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder(sb.DESC))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 versioning records, got %d", len(list))
	}
	if list[0].Content() != content2 {
		t.Errorf("expected content %q, got %q", content2, list[0].Content())
	}
}

func TestVersioningTrackEntity_VersioningDisabled(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   false,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()
	post := NewPost()
	post.SetTitle("Test Title").SetContent("Test Content")

	err = s.versioningTrackEntity(ctx, VERSIONING_TYPE_POST, "post-789", post)
	if err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestVersioningTrackEntity_Success(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()
	entityID := "post-track-001"
	post := NewPost()
	post.SetTitle("Track Test").SetContent("Track Content")

	err = s.versioningTrackEntity(ctx, VERSIONING_TYPE_POST, entityID, post)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := s.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_POST).
		SetEntityID(entityID))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 versioning record, got %d", len(list))
	}
}

func TestVersioningTrackEntity_UnsupportedEntity(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()
	type unsupportedEntity struct{}

	err = s.versioningTrackEntity(ctx, VERSIONING_TYPE_POST, "post-999", &unsupportedEntity{})
	if err == nil {
		t.Error("expected error, got nil")
	} else if !strings.Contains(err.Error(), "entity does not support versioning") {
		t.Errorf("expected error to contain 'entity does not support versioning', got %q", err.Error())
	}
}

func TestVersioningCreate(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	version := NewVersioning().
		SetEntityID("post-001").
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(`{"title":"Test"}`)

	err = store.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if version.ID() == "" {
		t.Error("expected version ID to be non-empty")
	}

	found, err := store.VersioningFindByID(ctx, version.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found == nil {
		t.Fatal("expected found to be non-nil")
	}
	if found.Content() != version.Content() {
		t.Errorf("expected content %q, got %q", version.Content(), found.Content())
	}
}

func TestVersioningCreate_NilStore(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	s.versioningStore = nil

	ctx := context.Background()
	version := NewVersioning().
		SetEntityID("post-001").
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(`{"title":"Test"}`)

	err = s.VersioningCreate(ctx, version)
	if err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestVersioningDelete(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	version := NewVersioning().
		SetEntityID("post-delete").
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(`{"title":"Delete Me"}`)

	err = store.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.VersioningDelete(ctx, version)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	found, err := store.VersioningFindByID(ctx, version.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found != nil {
		t.Error("expected found to be nil")
	}
}

func TestVersioningDelete_NilStore(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()

	version := NewVersioning().
		SetEntityID("post-delete-nil").
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(`{"title":"Delete Me"}`)

	err = s.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s.versioningStore = nil

	err = s.VersioningDelete(ctx, version)
	if err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestVersioningDeleteByID(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	version := NewVersioning().
		SetEntityID("post-delete-id").
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(`{"title":"Delete By ID"}`)

	err = store.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	versionID := version.ID()

	err = store.VersioningDeleteByID(ctx, versionID)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	found, err := store.VersioningFindByID(ctx, versionID)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found != nil {
		t.Error("expected found to be nil")
	}
}

func TestVersioningSoftDelete(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()

	version := NewVersioning().
		SetEntityID("post-soft-delete").
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(`{"title":"Soft Delete Me"}`)

	err = s.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = s.VersioningSoftDelete(ctx, version)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if version.SoftDeletedAt() == "" {
		t.Error("expected SoftDeletedAt to be non-empty")
	}
}

func TestVersioningSoftDelete_NilStore(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()

	version := NewVersioning().
		SetEntityID("post-soft-delete-nil").
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(`{"title":"Soft Delete Nil Store"}`)

	err = s.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s.versioningStore = nil

	err = s.VersioningSoftDelete(ctx, version)
	if err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestVersioningSoftDeleteByID(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()

	version := NewVersioning().
		SetEntityID("post-soft-delete-id").
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(`{"title":"Soft Delete By ID"}`)

	err = s.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	versionID := version.ID()

	err = s.VersioningSoftDeleteByID(ctx, versionID)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	found, err := s.VersioningFindByID(ctx, versionID)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found == nil {
		t.Fatal("expected found to be non-nil")
	}
	if found.SoftDeletedAt() == "" {
		t.Error("expected SoftDeletedAt to be non-empty")
	}
}

func TestVersioningSoftDeleteByID_NilStore(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()

	version := NewVersioning().
		SetEntityID("post-soft-delete-id-nil").
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(`{"title":"Soft Delete By ID"}`)

	err = s.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	versionID := version.ID()

	s.versioningStore = nil

	err = s.VersioningSoftDeleteByID(ctx, versionID)
	if err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestVersioningUpdate(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	version := NewVersioning().
		SetEntityID("post-update").
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(`{"title":"Original Title"}`)

	err = store.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	version.SetContent(`{"title":"Updated Title"}`)
	err = store.VersioningUpdate(ctx, version)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	found, err := store.VersioningFindByID(ctx, version.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found.Content() != `{"title":"Updated Title"}` {
		t.Errorf("expected content '{\"title\":\"Updated Title\"}', got %q", found.Content())
	}
}

func TestVersioningUpdate_NilStore(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()

	version := NewVersioning().
		SetEntityID("post-update-nil").
		SetEntityType(VERSIONING_TYPE_POST).
		SetContent(`{"title":"Original Title"}`)

	err = s.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s.versioningStore = nil

	err = s.VersioningUpdate(ctx, version)
	if err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestVersioningCreateIfChanged_NoExistingVersions(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	s, ok := store.(*storeImplementation)
	if !ok {
		t.Fatal("store is not *storeImplementation")
	}

	ctx := context.Background()
	entityID := "post-new-version"
	content := `{"title":"First Version"}`

	err = s.versioningCreateIfChanged(ctx, VERSIONING_TYPE_POST, entityID, content)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := s.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_POST).
		SetEntityID(entityID))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 versioning record, got %d", len(list))
	}
	if list[0].Content() != content {
		t.Errorf("expected content %q, got %q", content, list[0].Content())
	}
}

func TestVersioningList_EmptyResult(t *testing.T) {
	db := initDB()
	defer db.Close()
	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_versioning",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	list, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_POST).
		SetEntityID("non-existent"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d items", len(list))
	}
}
