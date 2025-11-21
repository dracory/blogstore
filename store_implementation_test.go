package blogstore

import (
	"context"
	"database/sql"
	"testing"

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

	postFound, errFind := store.PostFindByID(context.Background(), post.ID())
	if errFind != nil {
		t.Error("unexpected error:", errFind)
	}
	if postFound == nil {
		t.Error("Post MUST NOT be nil")
	}

	if postFound == nil {
		t.Error("BlogPost MUST NOT be nil")
	}

	if postFound.Title() != "2nd article" {
		t.Error("BlogPost first name MUST BE 'John', found: ", postFound.Title())
	}

	if postFound.Summary() != "Post Summary" {
		t.Error("BlogPost summary MUST BE 'Post Summary', found: ", postFound.Summary())
	}

	if postFound.Content() != "Post Content" {
		t.Error("BlogPost content MUST BE 'Post Content', found: ", postFound.Content())
	}

	if postFound.AuthorID() != "Post Author ID" {
		t.Error("BlogPost author ID MUST BE 'Post Content', found: ", postFound.AuthorID())
	}

	if postFound.ImageUrl() != "http://test.com/test.png" {
		t.Error("BlogPost image URL MUST BE 'http://test.com/test.png', found: ", postFound.ImageUrl())
	}

	if postFound.Status() != POST_STATUS_UNPUBLISHED {
		t.Error("BlogPost status MUST BE 'Unpublished', found: ", postFound.Status())
	}

	if postFound.CreatedAt() == "" {
		t.Error("BlogPost created MUST NOT BE empty, found: ", postFound.CreatedAt())
	}

	if postFound.UpdatedAt() == "" {
		t.Error("BlogPost updated MUST NOT BE empty, found: ", postFound.UpdatedAt())
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
	posts := []*Post{
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
	if list[0].Title() != "Post 1" {
		t.Errorf("PostList()[0].Title() = %q, want %q", list[0].Title(), "Post 1")
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

	found, err := store.PostFindByID(ctx, post.ID())
	if err != nil {
		t.Fatalf("PostFindByID() error = %v, want nil", err)
	}
	if found == nil {
		t.Fatalf("PostFindByID() returned nil, want non-nil")
	}
	if found.Status() != POST_STATUS_TRASH {
		t.Errorf("Status after PostTrash() = %q, want %q", found.Status(), POST_STATUS_TRASH)
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

	if err := store.PostSoftDeleteByID(ctx, post.ID()); err != nil {
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

	found, err := store.PostFindByID(ctx, post.ID())
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

	if err := store.PostDeleteByID(ctx, post2.ID()); err != nil {
		t.Fatalf("PostDeleteByID() error = %v, want nil", err)
	}

	found2, err := store.PostFindByID(ctx, post2.ID())
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

	// need concrete *store to call navigation helpers
	store, ok := storeIface.(*store)
	if !ok {
		t.Fatalf("NewStore did not return *store concrete type")
	}

	ctx := context.Background()

	// create three posts, then update CreatedAt so ordering is deterministic
	p1 := NewPost().SetTitle("First")
	p2 := NewPost().SetTitle("Second")
	p3 := NewPost().SetTitle("Third")

	for _, p := range []*Post{p1, p2, p3} {
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
	mid, err := store.PostFindByID(ctx, p2.ID())
	if err != nil {
		t.Fatalf("PostFindByID() error = %v, want nil", err)
	}
	if mid == nil {
		t.Fatalf("PostFindByID() returned nil for middle post")
	}

	prev, err := store.PostFindPrevious(*mid)
	if err != nil {
		t.Fatalf("PostFindPrevious() error = %v, want nil", err)
	}
	if prev == nil {
		t.Fatalf("PostFindPrevious() returned nil, want previous post")
	}
	if prev.Title() != "First" {
		t.Errorf("PostFindPrevious() Title = %q, want %q", prev.Title(), "First")
	}

	next, err := store.PostFindNext(*mid)
	if err != nil {
		t.Fatalf("PostFindNext() error = %v, want nil", err)
	}
	if next == nil {
		t.Fatalf("PostFindNext() returned nil, want next post")
	}
	if next.Title() != "Third" {
		t.Errorf("PostFindNext() Title = %q, want %q", next.Title(), "Third")
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

	for _, p := range []*Post{p1, p2, p3} {
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
	if list[0].Title() != "Golang testing" || list[2].Title() != "Irrelevant" {
		t.Errorf("PostList() ordering unexpected titles: first=%q last=%q", list[0].Title(), list[2].Title())
	}

	// soft delete one and verify WithDeleted behaviour
	if err := store.PostSoftDeleteByID(ctx, p1.ID()); err != nil {
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
