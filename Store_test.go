package blogstore

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestBlogRepositoryBlogPostCreate(t *testing.T) {
	db := initDB(":memory:")

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

	err = store.PostCreate(post)

	if err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestBlogRepositoryBlogPostFindByID(t *testing.T) {
	db := initDB(":memory:")

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

	err = store.PostCreate(post)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	postFound, errFind := store.PostFindByID(post.ID())
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
