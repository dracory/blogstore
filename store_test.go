package blogstore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/sb"
	"github.com/dracory/versionstore"
	_ "modernc.org/sqlite"
)

func initDB() *sql.DB {
	dsn := ":memory:" + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreVersioningCreateAndList(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:       "blog_posts",
		VersioningTableName: "blog_posts_version",
		VersioningEnabled:   true,
		DB:                  db,
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !store.VersioningEnabled() {
		t.Fatal("expected versioning to be enabled")
	}

	post := NewPost().
		SetTitle("Versioned post").
		SetStatus(POST_STATUS_DRAFT)

	if err := store.PostCreate(context.Background(), post); err != nil {
		t.Fatal("unexpected error:", err)
	}

	content, err := post.MarshalToVersioning()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.VersioningList(context.Background(), NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_POST).
		SetEntityID(post.GetID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder(sb.DESC))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 versioning record, got %d", len(list))
	}

	if list[0].Content() != content {
		t.Fatal("unexpected versioning content")
	}
}

func TestBlogRepositoryBlogPostCreate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	post := NewPost().
		SetStatus(POST_STATUS_UNPUBLISHED).
		SetTitle("1st article")

	err = store.PostCreate(context.Background(), post)

	if err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestBlogRepositoryBlogPostFindByID(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	post := NewPost().
		SetStatus(POST_STATUS_UNPUBLISHED).
		SetTitle("2nd article").
		SetContent("Post Content").
		SetSummary("Post Summary").
		SetAuthorID("Post Author ID").
		SetImageUrl("http://test.com/test.png").
		SetFeatured(YES)

	err = store.PostCreate(context.Background(), post)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	postFound, errFind := store.PostFindByID(context.Background(), post.GetID())
	if errFind != nil {
		t.Error("unexpected error:", errFind)
	}
	if postFound == nil {
		t.Error("Post MUST NOT be nil")
	}

	if postFound == nil {
		t.Error("BlogPost MUST NOT be nil")
	}

	if postFound.GetTitle() != "2nd article" {
		t.Error("BlogPost first name MUST BE 'John', found: ", postFound.GetTitle())
	}

	if postFound.GetSummary() != "Post Summary" {
		t.Error("BlogPost summary MUST BE 'Post Summary', found: ", postFound.GetSummary())
	}

	if postFound.GetContent() != "Post Content" {
		t.Error("BlogPost content MUST BE 'Post Content', found: ", postFound.GetContent())
	}

	if postFound.GetAuthorID() != "Post Author ID" {
		t.Error("BlogPost author ID MUST BE 'Post Content', found: ", postFound.GetAuthorID())
	}

	if postFound.GetImageUrl() != "http://test.com/test.png" {
		t.Error("BlogPost image URL MUST BE 'http://test.com/test.png', found: ", postFound.GetImageUrl())
	}

	if postFound.GetStatus() != POST_STATUS_UNPUBLISHED {
		t.Error("BlogPost status MUST BE 'Unpublished', found: ", postFound.GetStatus())
	}

	if postFound.GetCreatedAt() == "" {
		t.Error("BlogPost created MUST NOT BE empty, found: ", postFound.GetCreatedAt())
	}

	if postFound.GetUpdatedAt() == "" {
		t.Error("BlogPost updated MUST NOT BE empty, found: ", postFound.GetUpdatedAt())
	}
}

func TestBlogRepositoryBlogPostFindByOldSlug(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	post := NewPost().
		SetStatus(POST_STATUS_PUBLISHED).
		SetTitle("My Test Post").
		SetContent("Post Content").
		SetAuthorID("Post Author ID")

	// Add old slugs
	if err := post.AddOldSlug("old-slug-1"); err != nil {
		t.Fatal("AddOldSlug() error:", err)
	}
	if err := post.AddOldSlug("old-slug-2"); err != nil {
		t.Fatal("AddOldSlug() error:", err)
	}

	err = store.PostCreate(context.Background(), post)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	// Find by old slug
	postFound, errFind := store.PostFindByOldSlug(context.Background(), "old-slug-1")
	if errFind != nil {
		t.Error("unexpected error:", errFind)
	}
	if postFound == nil {
		t.Error("Post MUST NOT be nil")
	}

	if postFound.GetTitle() != "My Test Post" {
		t.Error("Post title MUST BE 'My Test Post', found: ", postFound.GetTitle())
	}

	// Test non-existent old slug
	notFound, errNotFound := store.PostFindByOldSlug(context.Background(), "non-existent-slug")
	if errNotFound != nil {
		t.Error("unexpected error for non-existent old slug:", errNotFound)
	}
	if notFound != nil {
		t.Error("Post SHOULD be nil for non-existent old slug")
	}
}

func TestBlogRepositoryBlogPostFindBySlug(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	post := NewPost().
		SetStatus(POST_STATUS_PUBLISHED).
		SetTitle("My Test Post").
		SetContent("Post Content").
		SetAuthorID("Post Author ID").
		SetSlug("my-test-post")

	err = store.PostCreate(context.Background(), post)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	// Find by slug
	postFound, errFind := store.PostFindBySlug(context.Background(), "my-test-post")
	if errFind != nil {
		t.Error("unexpected error:", errFind)
	}
	if postFound == nil {
		t.Error("Post MUST NOT be nil")
	}

	if postFound.GetTitle() != "My Test Post" {
		t.Error("Post title MUST BE 'My Test Post', found: ", postFound.GetTitle())
	}

	if postFound.GetSlug() != "my-test-post" {
		t.Error("Post slug MUST BE 'my-test-post', found: ", postFound.GetSlug())
	}

	// Test non-existent slug
	notFound, errNotFound := store.PostFindBySlug(context.Background(), "non-existent-slug")
	if errNotFound != nil {
		t.Error("unexpected error for non-existent slug:", errNotFound)
	}
	if notFound != nil {
		t.Error("Post SHOULD be nil for non-existent slug")
	}

	// Test empty slug
	emptySlug, errEmptySlug := store.PostFindBySlug(context.Background(), "")
	if errEmptySlug == nil {
		t.Error("expected error for empty slug, got nil")
	}
	if emptySlug != nil {
		t.Error("Post SHOULD be nil for empty slug")
	}
}

func TestStorePostListAndCount(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// create multiple posts with different statuses
	posts := []PostInterface{
		NewPost().SetTitle("Post 1").SetStatus(POST_STATUS_PUBLISHED),
		NewPost().SetTitle("Post 2").SetStatus(POST_STATUS_UNPUBLISHED),
		NewPost().SetTitle("Post 3").SetStatus(POST_STATUS_DRAFT),
	}

	for _, p := range posts {
		if err := store.PostCreate(ctx, p); err != nil {
			t.Fatalf("PostCreate() error = %v, want nil", err)
		}
	}

	// total count
	count, err := store.PostCount(ctx, PostQueryOptions{})
	if err != nil {
		t.Fatalf("PostCount() error = %v, want nil", err)
	}
	if count != 3 {
		t.Fatalf("PostCount() = %d, want %d", count, 3)
	}

	// filter by status
	list, err := store.PostList(ctx, PostQueryOptions{Status: POST_STATUS_PUBLISHED})
	if err != nil {
		t.Fatalf("PostList() error = %v, want nil", err)
	}
	if len(list) != 1 {
		t.Fatalf("PostList() len = %d, want %d", len(list), 1)
	}
	if list[0].GetTitle() != "Post 1" {
		t.Errorf("PostList()[0].Title() = %q, want %q", list[0].GetTitle(), "Post 1")
	}
}

func TestStorePostListMetaEquals(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// create posts with different metadata
	post1 := NewPost().SetTitle("Post with wp_id").SetStatus(POST_STATUS_PUBLISHED)
	if err := post1.SetMetas(map[string]string{"wp_id": "123", "wp_post_name": "my-post"}); err != nil {
		t.Fatalf("SetMetas() error = %v", err)
	}
	if err := store.PostCreate(ctx, post1); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	post2 := NewPost().SetTitle("Post with different wp_id").SetStatus(POST_STATUS_PUBLISHED)
	if err := post2.SetMetas(map[string]string{"wp_id": "456", "custom_field": "value"}); err != nil {
		t.Fatalf("SetMetas() error = %v", err)
	}
	if err := store.PostCreate(ctx, post2); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	post3 := NewPost().SetTitle("Post without metadata").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post3); err != nil {
		t.Fatalf("PostCreate() error = %v", err)
	}

	// Test MetaEquals with wp_id
	list, err := store.PostList(ctx, PostQueryOptions{
		MetaEquals: map[string]string{"wp_id": "123"},
	})
	if err != nil {
		t.Fatalf("PostList() with MetaEquals error = %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("PostList() with MetaEquals len = %d, want 1", len(list))
	}
	if list[0].GetID() != post1.GetID() {
		t.Errorf("PostList() with MetaEquals returned wrong post, want ID %s, got %s", post1.GetID(), list[0].GetID())
	}

	// Test MetaEquals with wp_post_name
	list, err = store.PostList(ctx, PostQueryOptions{
		MetaEquals: map[string]string{"wp_post_name": "my-post"},
	})
	if err != nil {
		t.Fatalf("PostList() with MetaEquals error = %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("PostList() with MetaEquals len = %d, want 1", len(list))
	}
	if list[0].GetID() != post1.GetID() {
		t.Errorf("PostList() with MetaEquals returned wrong post, want ID %s, got %s", post1.GetID(), list[0].GetID())
	}

	// Test MetaEquals with non-existent value
	list, err = store.PostList(ctx, PostQueryOptions{
		MetaEquals: map[string]string{"wp_id": "999"},
	})
	if err != nil {
		t.Fatalf("PostList() with MetaEquals error = %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("PostList() with non-existent MetaEquals len = %d, want 0", len(list))
	}
}

func TestStorePostTrashAndUpdate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().
		SetStatus(POST_STATUS_PUBLISHED).
		SetTitle("Trash Me")

	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}

	if err := store.PostTrash(ctx, post); err != nil {
		t.Fatalf("PostTrash() error = %v, want nil", err)
	}

	found, err := store.PostFindByID(ctx, post.GetID())
	if err != nil {
		t.Fatalf("PostFindByID() error = %v, want nil", err)
	}
	if found == nil {
		t.Fatalf("PostFindByID() returned nil, want non-nil")
	}
	if found.GetStatus() != POST_STATUS_TRASH {
		t.Errorf("Status after PostTrash() = %q, want %q", found.GetStatus(), POST_STATUS_TRASH)
	}
}

func TestStorePostSoftDeleteAndDelete(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	post := NewPost().
		SetStatus(POST_STATUS_PUBLISHED).
		SetTitle("Soft Delete Me")

	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}

	// ensure it is visible before soft delete
	count, err := store.PostCount(ctx, PostQueryOptions{})
	if err != nil {
		t.Fatalf("PostCount() error = %v, want nil", err)
	}
	if count != 1 {
		t.Fatalf("PostCount() before soft delete = %d, want %d", count, 1)
	}

	if err := store.PostSoftDeleteByID(ctx, post.GetID()); err != nil {
		t.Fatalf("PostSoftDeleteByID() error = %v, want nil", err)
	}

	// soft-deleted posts should not be returned by default queries
	count, err = store.PostCount(ctx, PostQueryOptions{})
	if err != nil {
		t.Fatalf("PostCount() after soft delete error = %v, want nil", err)
	}
	if count != 0 {
		t.Fatalf("PostCount() after soft delete = %d, want %d", count, 0)
	}

	found, err := store.PostFindByID(ctx, post.GetID())
	if err != nil {
		t.Fatalf("PostFindByID() after soft delete error = %v, want nil", err)
	}
	if found != nil {
		t.Fatalf("PostFindByID() after soft delete = %#v, want nil", found)
	}

	// now create a new post and delete by ID
	post2 := NewPost().
		SetStatus(POST_STATUS_PUBLISHED).
		SetTitle("Delete Me")

	if err := store.PostCreate(ctx, post2); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}

	if err := store.PostDeleteByID(ctx, post2.GetID()); err != nil {
		t.Fatalf("PostDeleteByID() error = %v, want nil", err)
	}

	found2, err := store.PostFindByID(ctx, post2.GetID())
	if err != nil {
		t.Fatalf("PostFindByID() after delete error = %v, want nil", err)
	}
	if found2 != nil {
		t.Fatalf("PostFindByID() after delete = %#v, want nil", found2)
	}
}

func TestStorePostFindPreviousAndNext(t *testing.T) {
	db := initDB()

	storeIface, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// need concrete *storeImplementation to call navigation helpers
	store, ok := storeIface.(*storeImplementation)
	if !ok {
		t.Fatalf("NewStore did not return *storeImplementation concrete type")
	}

	ctx := context.Background()

	// create three posts, then update CreatedAt so ordering is deterministic
	p1 := NewPost().SetTitle("First")
	p2 := NewPost().SetTitle("Second")
	p3 := NewPost().SetTitle("Third")

	for _, p := range []PostInterface{p1, p2, p3} {
		if err := store.PostCreate(ctx, p); err != nil {
			t.Fatalf("PostCreate() error = %v, want nil", err)
		}
	}

	// now set explicit CreatedAt values and persist via PostUpdate
	p1.SetCreatedAt("2020-01-01 00:00:00")
	if err := store.PostUpdate(ctx, p1); err != nil {
		t.Fatalf("PostUpdate() for p1 error = %v, want nil", err)
	}
	p2.SetCreatedAt("2020-01-02 00:00:00")
	if err := store.PostUpdate(ctx, p2); err != nil {
		t.Fatalf("PostUpdate() for p2 error = %v, want nil", err)
	}
	p3.SetCreatedAt("2020-01-03 00:00:00")
	if err := store.PostUpdate(ctx, p3); err != nil {
		t.Fatalf("PostUpdate() for p3 error = %v, want nil", err)
	}

	// reload middle post from DB to ensure timestamps persisted
	mid, err := store.PostFindByID(ctx, p2.GetID())
	if err != nil {
		t.Fatalf("PostFindByID() error = %v, want nil", err)
	}
	if mid == nil {
		t.Fatalf("PostFindByID() returned nil for middle post")
	}

	prev, err := store.PostFindPrevious(mid)
	if err != nil {
		t.Fatalf("PostFindPrevious() error = %v, want nil", err)
	}
	if prev == nil {
		t.Fatalf("PostFindPrevious() returned nil, want previous post")
	}
	if prev.GetTitle() != "First" {
		t.Errorf("PostFindPrevious() Title = %q, want %q", prev.GetTitle(), "First")
	}

	next, err := store.PostFindNext(mid)
	if err != nil {
		t.Fatalf("PostFindNext() error = %v, want nil", err)
	}
	if next == nil {
		t.Fatalf("PostFindNext() returned nil, want next post")
	}
	if next.GetTitle() != "Third" {
		t.Errorf("PostFindNext() Title = %q, want %q", next.GetTitle(), "Third")
	}
}

func TestStorePostListSearchOrderingAndWithDeleted(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	p1 := NewPost().
		SetTitle("Golang testing").
		SetContent("Content about go")
	p2 := NewPost().
		SetTitle("Another post").
		SetContent("Search me please")
	p3 := NewPost().
		SetTitle("Irrelevant").
		SetContent("Nothing to see here")

	for _, p := range []PostInterface{p1, p2, p3} {
		if err := store.PostCreate(ctx, p); err != nil {
			t.Fatalf("PostCreate() error = %v, want nil", err)
		}
	}

	// set explicit CreatedAt values and verify ordering
	p1.SetCreatedAt("2020-01-01 10:00:00")
	if err := store.PostUpdate(ctx, p1); err != nil {
		t.Fatalf("PostUpdate() for p1 error = %v, want nil", err)
	}
	p2.SetCreatedAt("2020-01-01 11:00:00")
	if err := store.PostUpdate(ctx, p2); err != nil {
		t.Fatalf("PostUpdate() for p2 error = %v, want nil", err)
	}
	p3.SetCreatedAt("2020-01-01 12:00:00")
	if err := store.PostUpdate(ctx, p3); err != nil {
		t.Fatalf("PostUpdate() for p3 error = %v, want nil", err)
	}

	list, err := store.PostList(ctx, PostQueryOptions{
		OrderBy:   COLUMN_CREATED_AT,
		SortOrder: "asc",
	})
	if err != nil {
		t.Fatalf("PostList() with ordering error = %v, want nil", err)
	}
	if len(list) != 3 {
		t.Fatalf("PostList() with ordering len = %d, want %d", len(list), 3)
	}
	if list[0].GetTitle() != "Golang testing" || list[2].GetTitle() != "Irrelevant" {
		t.Errorf("PostList() ordering unexpected titles: first=%q last=%q", list[0].GetTitle(), list[2].GetTitle())
	}

	// soft delete one and verify WithDeleted behaviour
	if err := store.PostSoftDeleteByID(ctx, p1.GetID()); err != nil {
		t.Fatalf("PostSoftDeleteByID() error = %v, want nil", err)
	}

	listDefault, err := store.PostList(ctx, PostQueryOptions{})
	if err != nil {
		t.Fatalf("PostList() default error = %v, want nil", err)
	}
	if len(listDefault) != 2 { // soft-deleted excluded
		t.Fatalf("PostList() default len = %d, want %d", len(listDefault), 2)
	}

	listWithDeleted, err := store.PostList(ctx, PostQueryOptions{WithDeleted: true})
	if err != nil {
		t.Fatalf("PostList() WithDeleted error = %v, want nil", err)
	}
	if len(listWithDeleted) != 3 {
		t.Fatalf("PostList() WithDeleted len = %d, want %d", len(listWithDeleted), 3)
	}
}
