package blogstore

import (
	"context"
	"testing"
)

func TestStoreMediaCreate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	media := NewMedia().
		SetEntityID(post.GetID()).
		SetTitle("test.jpg").
		SetURL("https://example.com/test.jpg").
		SetType("image/jpeg").
		SetSize("1024").
		SetExtension("jpg").
		SetSequence(1)

	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatalf("MediaCreate() error = %v", err)
	}

	if media.GetID() == "" {
		t.Error("MediaCreate() should set ID")
	}
	if media.GetCreatedAt() == "" {
		t.Error("MediaCreate() should set CreatedAt")
	}
	if media.GetUpdatedAt() == "" {
		t.Error("MediaCreate() should set UpdatedAt")
	}
}

func TestStoreMediaCreateErrors(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Test with empty entity_id
	media := NewMedia().SetEntityID("").SetTitle("test.jpg").SetURL("https://example.com/test.jpg")
	if err := store.MediaCreate(ctx, media); err == nil {
		t.Error("MediaCreate() with empty entity_id should return error")
	}

	// Test with nil
	if err := store.MediaCreate(ctx, nil); err == nil {
		t.Error("MediaCreate() with nil should return error")
	}
}

func TestStoreMediaFindByID(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	media := NewMedia().
		SetEntityID(post.GetID()).
		SetTitle("find.jpg").
		SetURL("https://example.com/find.jpg").
		SetType("image/jpeg").
		SetExtension("jpg").
		SetSequence(1)

	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatalf("MediaCreate() error = %v", err)
	}

	found, err := store.MediaFindByID(ctx, media.GetID())
	if err != nil {
		t.Fatalf("MediaFindByID() error = %v", err)
	}
	if found == nil {
		t.Fatal("MediaFindByID() returned nil")
	}
	if found.GetTitle() != "find.jpg" {
		t.Errorf("GetTitle() = %q, want %q", found.GetTitle(), "find.jpg")
	}
	if found.GetEntityID() != post.GetID() {
		t.Errorf("GetEntityID() = %q, want %q", found.GetEntityID(), post.GetID())
	}
	if found.GetURL() != "https://example.com/find.jpg" {
		t.Errorf("GetURL() = %q, want %q", found.GetURL(), "https://example.com/find.jpg")
	}
}

func TestStoreMediaList(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	// Create multiple media
	for i := 0; i < 3; i++ {
		m := NewMedia().
			SetEntityID(post.GetID()).
			SetTitle("file-" + string(rune('A'+i)) + ".jpg").
			SetURL("https://example.com/file" + string(rune('A'+i)) + ".jpg").
			SetExtension("jpg").
			SetSequence(i + 1)
		if err := store.MediaCreate(ctx, m); err != nil {
			t.Fatalf("MediaCreate() error = %v", err)
		}
	}

	// List all
	list, err := store.MediaList(ctx, MediaQueryOptions{})
	if err != nil {
		t.Fatalf("MediaList() error = %v", err)
	}
	if len(list) != 3 {
		t.Errorf("MediaList() len = %d, want 3", len(list))
	}

	// List by entity ID
	list, err = store.MediaListByEntityID(ctx, post.GetID())
	if err != nil {
		t.Fatalf("MediaListByEntityID() error = %v", err)
	}
	if len(list) != 3 {
		t.Errorf("MediaListByEntityID() len = %d, want 3", len(list))
	}
}

func TestStoreMediaCount(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	// Create 3 media
	for i := 0; i < 3; i++ {
		m := NewMedia().
			SetEntityID(post.GetID()).
			SetTitle("count-" + string(rune('A'+i)) + ".jpg").
			SetURL("https://example.com/count.jpg").
			SetExtension("jpg").
			SetSequence(i + 1)
		if err := store.MediaCreate(ctx, m); err != nil {
			t.Fatalf("MediaCreate() error = %v", err)
		}
	}

	count, err := store.MediaCount(ctx, MediaQueryOptions{})
	if err != nil {
		t.Fatalf("MediaCount() error = %v", err)
	}
	if count != 3 {
		t.Errorf("MediaCount() = %d, want 3", count)
	}

	// Count by entity ID
	count, err = store.MediaCount(ctx, MediaQueryOptions{EntityID: post.GetID()})
	if err != nil {
		t.Fatalf("MediaCount() error = %v", err)
	}
	if count != 3 {
		t.Errorf("MediaCount() by entity = %d, want 3", count)
	}
}

func TestStoreMediaUpdate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	media := NewMedia().
		SetEntityID(post.GetID()).
		SetTitle("original.jpg").
		SetURL("https://example.com/original.jpg").
		SetExtension("jpg").
		SetSequence(1)
	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatalf("MediaCreate() error = %v", err)
	}

	media.SetTitle("updated.jpg")
	media.SetURL("https://example.com/updated.jpg")
	media.SetStatus(MEDIA_STATUS_ACTIVE)

	if err := store.MediaUpdate(ctx, media); err != nil {
		t.Fatalf("MediaUpdate() error = %v", err)
	}

	found, err := store.MediaFindByID(ctx, media.GetID())
	if err != nil {
		t.Fatalf("MediaFindByID() error = %v", err)
	}
	if found.GetTitle() != "updated.jpg" {
		t.Errorf("GetTitle() = %q, want %q", found.GetTitle(), "updated.jpg")
	}
	if found.GetURL() != "https://example.com/updated.jpg" {
		t.Errorf("GetURL() = %q, want %q", found.GetURL(), "https://example.com/updated.jpg")
	}
	if found.GetStatus() != MEDIA_STATUS_ACTIVE {
		t.Errorf("GetStatus() = %q, want %q", found.GetStatus(), MEDIA_STATUS_ACTIVE)
	}
}

func TestStoreMediaSoftDelete(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	media := NewMedia().
		SetEntityID(post.GetID()).
		SetTitle("softdelete.jpg").
		SetURL("https://example.com/softdelete.jpg").
		SetExtension("jpg").
		SetSequence(1)
	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatalf("MediaCreate() error = %v", err)
	}

	// Soft delete
	if err := store.MediaSoftDelete(ctx, media); err != nil {
		t.Fatalf("MediaSoftDelete() error = %v", err)
	}

	// Should not appear in normal list
	list, err := store.MediaList(ctx, MediaQueryOptions{})
	if err != nil {
		t.Fatalf("MediaList() error = %v", err)
	}
	if len(list) != 0 {
		t.Errorf("MediaList() after soft delete len = %d, want 0", len(list))
	}

	// Should appear with WithDeleted
	list, err = store.MediaList(ctx, MediaQueryOptions{WithDeleted: true})
	if err != nil {
		t.Fatalf("MediaList() with deleted error = %v", err)
	}
	if len(list) != 1 {
		t.Errorf("MediaList() with deleted len = %d, want 1", len(list))
	}
}

func TestStoreMediaSoftDeleteByID(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	media := NewMedia().
		SetEntityID(post.GetID()).
		SetTitle("softdeletebyid.jpg").
		SetURL("https://example.com/softdeletebyid.jpg").
		SetExtension("jpg").
		SetSequence(1)
	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatalf("MediaCreate() error = %v", err)
	}

	if err := store.MediaSoftDeleteByID(ctx, media.GetID()); err != nil {
		t.Fatalf("MediaSoftDeleteByID() error = %v", err)
	}

	list, err := store.MediaList(ctx, MediaQueryOptions{})
	if err != nil {
		t.Fatalf("MediaList() error = %v", err)
	}
	if len(list) != 0 {
		t.Errorf("MediaList() after soft delete by ID len = %d, want 0", len(list))
	}
}

func TestStoreMediaDelete(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	media := NewMedia().
		SetEntityID(post.GetID()).
		SetTitle("delete.jpg").
		SetURL("https://example.com/delete.jpg").
		SetExtension("jpg").
		SetSequence(1)
	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatalf("MediaCreate() error = %v", err)
	}

	// Hard delete
	if err := store.MediaDeleteByID(ctx, media.GetID()); err != nil {
		t.Fatalf("MediaDeleteByID() error = %v", err)
	}

	// Should not find it even with WithDeleted
	list, err := store.MediaList(ctx, MediaQueryOptions{WithDeleted: true})
	if err != nil {
		t.Fatalf("MediaList() error = %v", err)
	}
	if len(list) != 0 {
		t.Errorf("MediaList() after hard delete with WithDeleted len = %d, want 0", len(list))
	}
}

func TestStoreMediaSearch(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	mediaItems := []MediaInterface{
		NewMedia().SetEntityID(post.GetID()).SetTitle("vacation-photo.jpg").SetURL("https://example.com/v1.jpg").SetExtension("jpg").SetSequence(1),
		NewMedia().SetEntityID(post.GetID()).SetTitle("vacation-video.mp4").SetURL("https://example.com/v2.mp4").SetExtension("mp4").SetSequence(2),
		NewMedia().SetEntityID(post.GetID()).SetTitle("document.pdf").SetURL("https://example.com/doc.pdf").SetExtension("pdf").SetSequence(3),
	}

	for _, m := range mediaItems {
		if err := store.MediaCreate(ctx, m); err != nil {
			t.Fatalf("MediaCreate() error = %v", err)
		}
	}

	// Search for "vacation"
	results, err := store.MediaList(ctx, MediaQueryOptions{Search: "vacation"})
	if err != nil {
		t.Fatalf("MediaList() search error = %v", err)
	}
	if len(results) != 2 {
		t.Errorf("MediaList() search 'vacation' len = %d, want 2", len(results))
	}

	// Search for "document"
	results, err = store.MediaList(ctx, MediaQueryOptions{Search: "document"})
	if err != nil {
		t.Fatalf("MediaList() search error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("MediaList() search 'document' len = %d, want 1", len(results))
	}

	// Search for non-existent
	results, err = store.MediaList(ctx, MediaQueryOptions{Search: "nonexistent"})
	if err != nil {
		t.Fatalf("MediaList() search error = %v", err)
	}
	if len(results) != 0 {
		t.Errorf("MediaList() search 'nonexistent' len = %d, want 0", len(results))
	}
}

func TestStoreMediaListPagination(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	for i := 0; i < 5; i++ {
		m := NewMedia().
			SetEntityID(post.GetID()).
			SetTitle("file.jpg").
			SetURL("https://example.com/file.jpg").
			SetExtension("jpg").
			SetSequence(i + 1)
		if err := store.MediaCreate(ctx, m); err != nil {
			t.Fatalf("MediaCreate() error = %v", err)
		}
	}

	// Test limit
	list, err := store.MediaList(ctx, MediaQueryOptions{Limit: 2})
	if err != nil {
		t.Fatalf("MediaList() error = %v", err)
	}
	if len(list) != 2 {
		t.Errorf("MediaList() with Limit=2 len = %d, want 2", len(list))
	}

	// Test offset
	list, err = store.MediaList(ctx, MediaQueryOptions{Limit: 2, Offset: 2})
	if err != nil {
		t.Fatalf("MediaList() error = %v", err)
	}
	if len(list) != 2 {
		t.Errorf("MediaList() with Offset=2, Limit=2 len = %d, want 2", len(list))
	}

	// Total count should still be 5
	count, err := store.MediaCount(ctx, MediaQueryOptions{})
	if err != nil {
		t.Fatalf("MediaCount() error = %v", err)
	}
	if count != 5 {
		t.Errorf("MediaCount() = %d, want 5", count)
	}
}

func TestStoreMediaDefaultTableName(t *testing.T) {
	db := initDB()

	// Don't set MediaTableName - should default to "blog_media"
	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store.GetMediaTableName() != "blog_media" {
		t.Errorf("GetMediaTableName() = %q, want %q", store.GetMediaTableName(), "blog_media")
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	media := NewMedia().
		SetEntityID(post.GetID()).
		SetTitle("default-table.jpg").
		SetURL("https://example.com/default.jpg").
		SetExtension("jpg").
		SetSequence(1)

	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatalf("MediaCreate() error = %v", err)
	}

	found, err := store.MediaFindByID(ctx, media.GetID())
	if err != nil {
		t.Fatalf("MediaFindByID() error = %v", err)
	}
	if found == nil {
		t.Fatal("MediaFindByID() returned nil")
	}
	if found.GetTitle() != "default-table.jpg" {
		t.Errorf("GetTitle() = %q, want %q", found.GetTitle(), "default-table.jpg")
	}
}

func TestStoreMediaSetGetTableName(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "custom_media",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store.GetMediaTableName() != "custom_media" {
		t.Errorf("GetMediaTableName() = %q, want %q", store.GetMediaTableName(), "custom_media")
	}

	// Test setting table name after creation
	store.SetMediaTableName("renamed_media")
	if store.GetMediaTableName() != "renamed_media" {
		t.Errorf("GetMediaTableName() after SetMediaTableName() = %q, want %q", store.GetMediaTableName(), "renamed_media")
	}
}
