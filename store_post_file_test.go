package blogstore

import (
	"context"
	"testing"
)

// ============================ POST FILE STORE TESTS ============================

func TestStorePostFileCreateAndFind(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		PostFileTableName:  "blog_post_file",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create a post first
	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}

	// Create a post file
	file := NewPostFile().
		SetPostID(post.GetID()).
		SetName("photo.jpg").
		SetURL("https://example.com/uploads/photo.jpg").
		SetType("image/jpeg").
		SetSize("102400").
		SetExtension("jpg").
		SetSequence(1)

	if err := store.PostFileCreate(ctx, file); err != nil {
		t.Fatalf("PostFileCreate() error = %v, want nil", err)
	}

	// Find by ID
	found, err := store.PostFileFindByID(ctx, file.GetID())
	if err != nil {
		t.Fatalf("PostFileFindByID() error = %v, want nil", err)
	}
	if found == nil {
		t.Fatal("PostFileFindByID() returned nil, want non-nil")
	}
	if found.GetPostID() != post.GetID() {
		t.Errorf("GetPostID() = %q, want %q", found.GetPostID(), post.GetID())
	}
	if found.GetName() != "photo.jpg" {
		t.Errorf("GetName() = %q, want %q", found.GetName(), "photo.jpg")
	}
	if found.GetURL() != "https://example.com/uploads/photo.jpg" {
		t.Errorf("GetURL() = %q, want %q", found.GetURL(), "https://example.com/uploads/photo.jpg")
	}
	if found.GetType() != "image/jpeg" {
		t.Errorf("GetType() = %q, want %q", found.GetType(), "image/jpeg")
	}
	if found.GetSize() != "102400" {
		t.Errorf("GetSize() = %q, want %q", found.GetSize(), "102400")
	}
	if found.GetExtension() != "jpg" {
		t.Errorf("GetExtension() = %q, want %q", found.GetExtension(), "jpg")
	}
	if found.GetSequence() != 1 {
		t.Errorf("GetSequence() = %d, want %d", found.GetSequence(), 1)
	}
	if found.GetCreatedAt() == "" {
		t.Error("GetCreatedAt() should not be empty")
	}
	if found.GetUpdatedAt() == "" {
		t.Error("GetUpdatedAt() should not be empty")
	}
}

func TestStorePostFileListAndCount(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		PostFileTableName:  "blog_post_file",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create a post
	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}

	// Create multiple post files
	files := []PostFileInterface{
		NewPostFile().SetPostID(post.GetID()).SetName("file1.jpg").SetURL("https://example.com/f1.jpg").SetType("image/jpeg").SetExtension("jpg").SetSequence(1),
		NewPostFile().SetPostID(post.GetID()).SetName("file2.png").SetURL("https://example.com/f2.png").SetType("image/png").SetExtension("png").SetSequence(2),
		NewPostFile().SetPostID(post.GetID()).SetName("file3.pdf").SetURL("https://example.com/f3.pdf").SetType("application/pdf").SetExtension("pdf").SetSequence(3),
	}

	for _, f := range files {
		if err := store.PostFileCreate(ctx, f); err != nil {
			t.Fatalf("PostFileCreate() error = %v, want nil", err)
		}
	}

	// Count all
	count, err := store.PostFileCount(ctx, PostFileQueryOptions{})
	if err != nil {
		t.Fatalf("PostFileCount() error = %v, want nil", err)
	}
	if count != 3 {
		t.Errorf("PostFileCount() = %d, want %d", count, 3)
	}

	// List all
	list, err := store.PostFileList(ctx, PostFileQueryOptions{})
	if err != nil {
		t.Fatalf("PostFileList() error = %v, want nil", err)
	}
	if len(list) != 3 {
		t.Errorf("PostFileList() len = %d, want %d", len(list), 3)
	}

	// Filter by PostID
	listByPost, err := store.PostFileList(ctx, PostFileQueryOptions{PostID: post.GetID()})
	if err != nil {
		t.Fatalf("PostFileList() by PostID error = %v, want nil", err)
	}
	if len(listByPost) != 3 {
		t.Errorf("PostFileList() by PostID len = %d, want %d", len(listByPost), 3)
	}

	// Filter by extension
	listByExt, err := store.PostFileList(ctx, PostFileQueryOptions{Extension: "jpg"})
	if err != nil {
		t.Fatalf("PostFileList() by Extension error = %v, want nil", err)
	}
	if len(listByExt) != 1 {
		t.Errorf("PostFileList() by Extension len = %d, want %d", len(listByExt), 1)
	}
	if listByExt[0].GetName() != "file1.jpg" {
		t.Errorf("PostFileList() by Extension name = %q, want %q", listByExt[0].GetName(), "file1.jpg")
	}

	// Filter by type
	listByType, err := store.PostFileList(ctx, PostFileQueryOptions{Type: "image/png"})
	if err != nil {
		t.Fatalf("PostFileList() by Type error = %v, want nil", err)
	}
	if len(listByType) != 1 {
		t.Errorf("PostFileList() by Type len = %d, want %d", len(listByType), 1)
	}
}

func TestStorePostFileListByPostID(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		PostFileTableName:  "blog_post_file",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create two posts
	post1 := NewPost().SetTitle("Post 1").SetStatus(POST_STATUS_PUBLISHED)
	post2 := NewPost().SetTitle("Post 2").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post1); err != nil {
		t.Fatalf("PostCreate() post1 error = %v", err)
	}
	if err := store.PostCreate(ctx, post2); err != nil {
		t.Fatalf("PostCreate() post2 error = %v", err)
	}

	// Create files for post1
	f1 := NewPostFile().SetPostID(post1.GetID()).SetName("a.jpg").SetURL("https://example.com/a.jpg").SetExtension("jpg").SetSequence(2)
	f2 := NewPostFile().SetPostID(post1.GetID()).SetName("b.jpg").SetURL("https://example.com/b.jpg").SetExtension("jpg").SetSequence(1)
	if err := store.PostFileCreate(ctx, f1); err != nil {
		t.Fatalf("PostFileCreate() f1 error = %v", err)
	}
	if err := store.PostFileCreate(ctx, f2); err != nil {
		t.Fatalf("PostFileCreate() f2 error = %v", err)
	}

	// Create a file for post2
	f3 := NewPostFile().SetPostID(post2.GetID()).SetName("c.jpg").SetURL("https://example.com/c.jpg").SetExtension("jpg").SetSequence(1)
	if err := store.PostFileCreate(ctx, f3); err != nil {
		t.Fatalf("PostFileCreate() f3 error = %v", err)
	}

	// List files for post1 - should be ordered by sequence ASC
	files, err := store.PostFileListByPostID(ctx, post1.GetID())
	if err != nil {
		t.Fatalf("PostFileListByPostID() error = %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("PostFileListByPostID() len = %d, want %d", len(files), 2)
	}
	// f2 has sequence 1, f1 has sequence 2
	if files[0].GetName() != "b.jpg" {
		t.Errorf("PostFileListByPostID()[0].GetName() = %q, want %q", files[0].GetName(), "b.jpg")
	}
	if files[1].GetName() != "a.jpg" {
		t.Errorf("PostFileListByPostID()[1].GetName() = %q, want %q", files[1].GetName(), "a.jpg")
	}

	// List files for post2
	files2, err := store.PostFileListByPostID(ctx, post2.GetID())
	if err != nil {
		t.Fatalf("PostFileListByPostID() post2 error = %v", err)
	}
	if len(files2) != 1 {
		t.Errorf("PostFileListByPostID() post2 len = %d, want %d", len(files2), 1)
	}
}

func TestStorePostFileUpdate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		PostFileTableName:  "blog_post_file",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create post and file
	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	file := NewPostFile().
		SetPostID(post.GetID()).
		SetName("original.jpg").
		SetURL("https://example.com/original.jpg").
		SetType("image/jpeg").
		SetSize("1000").
		SetExtension("jpg").
		SetSequence(1)

	if err := store.PostFileCreate(ctx, file); err != nil {
		t.Fatalf("PostFileCreate() error = %v", err)
	}

	// Update
	file.SetName("updated.jpg").
		SetURL("https://example.com/updated.jpg").
		SetSize("2000").
		SetSequence(5)

	if err := store.PostFileUpdate(ctx, file); err != nil {
		t.Fatalf("PostFileUpdate() error = %v", err)
	}

	// Verify
	found, err := store.PostFileFindByID(ctx, file.GetID())
	if err != nil {
		t.Fatalf("PostFileFindByID() error = %v", err)
	}
	if found == nil {
		t.Fatal("PostFileFindByID() returned nil")
	}
	if found.GetName() != "updated.jpg" {
		t.Errorf("GetName() = %q, want %q", found.GetName(), "updated.jpg")
	}
	if found.GetURL() != "https://example.com/updated.jpg" {
		t.Errorf("GetURL() = %q, want %q", found.GetURL(), "https://example.com/updated.jpg")
	}
	if found.GetSize() != "2000" {
		t.Errorf("GetSize() = %q, want %q", found.GetSize(), "2000")
	}
	if found.GetSequence() != 5 {
		t.Errorf("GetSequence() = %d, want %d", found.GetSequence(), 5)
	}
}

func TestStorePostFileSoftDeleteAndDelete(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		PostFileTableName:  "blog_post_file",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create post and file
	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	file := NewPostFile().
		SetPostID(post.GetID()).
		SetName("delete-me.jpg").
		SetURL("https://example.com/delete-me.jpg").
		SetExtension("jpg").
		SetSequence(1)

	if err := store.PostFileCreate(ctx, file); err != nil {
		t.Fatalf("PostFileCreate() error = %v", err)
	}

	// Verify it exists
	count, err := store.PostFileCount(ctx, PostFileQueryOptions{})
	if err != nil {
		t.Fatalf("PostFileCount() error = %v", err)
	}
	if count != 1 {
		t.Fatalf("PostFileCount() before soft delete = %d, want 1", count)
	}

	// Soft delete by ID
	if err := store.PostFileSoftDeleteByID(ctx, file.GetID()); err != nil {
		t.Fatalf("PostFileSoftDeleteByID() error = %v", err)
	}

	// Should not be visible by default
	count, err = store.PostFileCount(ctx, PostFileQueryOptions{})
	if err != nil {
		t.Fatalf("PostFileCount() after soft delete error = %v", err)
	}
	if count != 0 {
		t.Fatalf("PostFileCount() after soft delete = %d, want 0", count)
	}

	// Should be visible with WithDeleted
	count, err = store.PostFileCount(ctx, PostFileQueryOptions{WithDeleted: true})
	if err != nil {
		t.Fatalf("PostFileCount() WithDeleted error = %v", err)
	}
	if count != 1 {
		t.Fatalf("PostFileCount() WithDeleted = %d, want 1", count)
	}

	// Permanently delete by ID
	if err := store.PostFileDeleteByID(ctx, file.GetID()); err != nil {
		t.Fatalf("PostFileDeleteByID() error = %v", err)
	}

	// Should not exist even with WithDeleted
	count, err = store.PostFileCount(ctx, PostFileQueryOptions{WithDeleted: true})
	if err != nil {
		t.Fatalf("PostFileCount() after delete error = %v", err)
	}
	if count != 0 {
		t.Fatalf("PostFileCount() after delete WithDeleted = %d, want 0", count)
	}
}

func TestStorePostFileSoftDeleteByObject(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		PostFileTableName:  "blog_post_file",
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

	file := NewPostFile().
		SetPostID(post.GetID()).
		SetName("soft-delete-obj.jpg").
		SetURL("https://example.com/soft-delete-obj.jpg").
		SetExtension("jpg").
		SetSequence(1)

	if err := store.PostFileCreate(ctx, file); err != nil {
		t.Fatalf("PostFileCreate() error = %v", err)
	}

	// Soft delete by object
	if err := store.PostFileSoftDelete(ctx, file); err != nil {
		t.Fatalf("PostFileSoftDelete() error = %v", err)
	}

	// Should not be visible by default
	found, err := store.PostFileFindByID(ctx, file.GetID())
	if err != nil {
		t.Fatalf("PostFileFindByID() error = %v", err)
	}
	if found != nil {
		t.Fatalf("PostFileFindByID() after soft delete = %#v, want nil", found)
	}

	// Should be visible with WithDeleted
	list, err := store.PostFileList(ctx, PostFileQueryOptions{WithDeleted: true})
	if err != nil {
		t.Fatalf("PostFileList() WithDeleted error = %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("PostFileList() WithDeleted len = %d, want 1", len(list))
	}
	if !list[0].IsSoftDeleted() {
		t.Error("IsSoftDeleted() should be true after soft delete")
	}
}

func TestStorePostFileDeleteByObject(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		PostFileTableName:  "blog_post_file",
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

	file := NewPostFile().
		SetPostID(post.GetID()).
		SetName("delete-obj.jpg").
		SetURL("https://example.com/delete-obj.jpg").
		SetExtension("jpg").
		SetSequence(1)

	if err := store.PostFileCreate(ctx, file); err != nil {
		t.Fatalf("PostFileCreate() error = %v", err)
	}

	// Delete by object
	if err := store.PostFileDelete(ctx, file); err != nil {
		t.Fatalf("PostFileDelete() error = %v", err)
	}

	// Should not exist
	found, err := store.PostFileFindByID(ctx, file.GetID())
	if err != nil {
		t.Fatalf("PostFileFindByID() error = %v", err)
	}
	if found != nil {
		t.Fatalf("PostFileFindByID() after delete = %#v, want nil", found)
	}
}

func TestStorePostFileCreateErrors(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		PostFileTableName:  "blog_post_file",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Test with empty post_id
	file := NewPostFile().SetPostID("").SetName("test.jpg").SetURL("https://example.com/test.jpg")
	if err := store.PostFileCreate(ctx, file); err == nil {
		t.Error("PostFileCreate() with empty post_id should return error")
	}

	// Test with nil
	if err := store.PostFileCreate(ctx, nil); err == nil {
		t.Error("PostFileCreate() with nil should return error")
	}
}

func TestStorePostFileSearch(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		PostFileTableName:  "blog_post_file",
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

	// Create files with different names
	files := []PostFileInterface{
		NewPostFile().SetPostID(post.GetID()).SetName("vacation-photo.jpg").SetURL("https://example.com/v1.jpg").SetExtension("jpg").SetSequence(1),
		NewPostFile().SetPostID(post.GetID()).SetName("vacation-video.mp4").SetURL("https://example.com/v2.mp4").SetExtension("mp4").SetSequence(2),
		NewPostFile().SetPostID(post.GetID()).SetName("document.pdf").SetURL("https://example.com/doc.pdf").SetExtension("pdf").SetSequence(3),
	}

	for _, f := range files {
		if err := store.PostFileCreate(ctx, f); err != nil {
			t.Fatalf("PostFileCreate() error = %v", err)
		}
	}

	// Search for "vacation"
	results, err := store.PostFileList(ctx, PostFileQueryOptions{Search: "vacation"})
	if err != nil {
		t.Fatalf("PostFileList() search error = %v", err)
	}
	if len(results) != 2 {
		t.Errorf("PostFileList() search 'vacation' len = %d, want 2", len(results))
	}

	// Search for "document"
	results, err = store.PostFileList(ctx, PostFileQueryOptions{Search: "document"})
	if err != nil {
		t.Fatalf("PostFileList() search error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("PostFileList() search 'document' len = %d, want 1", len(results))
	}

	// Search for non-existent
	results, err = store.PostFileList(ctx, PostFileQueryOptions{Search: "nonexistent"})
	if err != nil {
		t.Fatalf("PostFileList() search error = %v", err)
	}
	if len(results) != 0 {
		t.Errorf("PostFileList() search 'nonexistent' len = %d, want 0", len(results))
	}
}

func TestStorePostFileListPagination(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		PostFileTableName:  "blog_post_file",
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

	// Create 5 files
	for i := 0; i < 5; i++ {
		f := NewPostFile().
			SetPostID(post.GetID()).
			SetName("file.jpg").
			SetURL("https://example.com/file.jpg").
			SetExtension("jpg").
			SetSequence(i + 1)
		if err := store.PostFileCreate(ctx, f); err != nil {
			t.Fatalf("PostFileCreate() error = %v", err)
		}
	}

	// Test limit
	list, err := store.PostFileList(ctx, PostFileQueryOptions{Limit: 2})
	if err != nil {
		t.Fatalf("PostFileList() error = %v", err)
	}
	if len(list) != 2 {
		t.Errorf("PostFileList() with Limit=2 len = %d, want 2", len(list))
	}

	// Test offset
	list, err = store.PostFileList(ctx, PostFileQueryOptions{Limit: 2, Offset: 2})
	if err != nil {
		t.Fatalf("PostFileList() error = %v", err)
	}
	if len(list) != 2 {
		t.Errorf("PostFileList() with Offset=2, Limit=2 len = %d, want 2", len(list))
	}

	// Total count should still be 5
	count, err := store.PostFileCount(ctx, PostFileQueryOptions{})
	if err != nil {
		t.Fatalf("PostFileCount() error = %v", err)
	}
	if count != 5 {
		t.Errorf("PostFileCount() = %d, want 5", count)
	}
}

func TestStorePostFileDefaultTableName(t *testing.T) {
	db := initDB()

	// Don't set PostFileTableName - should default to "blog_post_file"
	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store.GetPostFileTableName() != "blog_post_file" {
		t.Errorf("GetPostFileTableName() = %q, want %q", store.GetPostFileTableName(), "blog_post_file")
	}

	ctx := context.Background()

	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	file := NewPostFile().
		SetPostID(post.GetID()).
		SetName("default-table.jpg").
		SetURL("https://example.com/default.jpg").
		SetExtension("jpg").
		SetSequence(1)

	if err := store.PostFileCreate(ctx, file); err != nil {
		t.Fatalf("PostFileCreate() error = %v", err)
	}

	found, err := store.PostFileFindByID(ctx, file.GetID())
	if err != nil {
		t.Fatalf("PostFileFindByID() error = %v", err)
	}
	if found == nil {
		t.Fatal("PostFileFindByID() returned nil")
	}
	if found.GetName() != "default-table.jpg" {
		t.Errorf("GetName() = %q, want %q", found.GetName(), "default-table.jpg")
	}
}

func TestStorePostFileSetGetTableName(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		PostFileTableName:  "custom_post_file",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store.GetPostFileTableName() != "custom_post_file" {
		t.Errorf("GetPostFileTableName() = %q, want %q", store.GetPostFileTableName(), "custom_post_file")
	}

	// Test setting table name after creation
	store.SetPostFileTableName("renamed_post_file")
	if store.GetPostFileTableName() != "renamed_post_file" {
		t.Errorf("GetPostFileTableName() after SetPostFileTableName() = %q, want %q", store.GetPostFileTableName(), "renamed_post_file")
	}
}
