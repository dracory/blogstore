package blogstore

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"
)

func TestExtractBlogMediaID(t *testing.T) {
	tests := []struct {
		name     string
		urlPath  string
		expected string
	}{
		{"short form png", "/blog/media/sbm2hdy7x10.png", "sbm2hdy7x10"},
		{"short form jpg", "/blog/media/abc123.jpg", "abc123"},
		{"short form no extension", "/blog/media/abc123", "abc123"},
		{"readable form", "/blog/media/sbm2hdy7x10/image.png", "sbm2hdy7x10"},
		{"readable form with underscores", "/blog/media/abc123/my_image_name.jpg", "abc123"},
		{"trailing slash", "/blog/media/sbm2hdy7x10/", "sbm2hdy7x10"},
		{"double slash", "/blog/media//sbm2hdy7x10.png", "sbm2hdy7x10"},
		{"empty path", "/blog/media/", ""},
		{"only extension", "/blog/media/.png", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractBlogMediaID(tt.urlPath)
			if got != tt.expected {
				t.Errorf("extractBlogMediaID(%q) = %q, want %q", tt.urlPath, got, tt.expected)
			}
		})
	}
}

func TestIsMediaURL(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/blog/media/abc123.png", true},
		{"/blog/media/abc123", true},
		{"/blog/media/abc123/handle.png", true},
		{"/blog/media/", true},
		{"/blog/post/123", false},
		{"/", false},
		{"/cms/media/abc123.png", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := IsMediaURL(tt.path)
			if got != tt.expected {
				t.Errorf("IsMediaURL(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestMediaServeURL(t *testing.T) {
	media := NewMedia()
	media.SetID("abc123")
	media.SetExtension(".png")

	got := media.ServeURL()
	expected := "/blog/media/abc123.png"
	if got != expected {
		t.Errorf("expected ServeURL %q, got %q", expected, got)
	}
}

func TestMediaServeURL_NoExtension(t *testing.T) {
	media := NewMedia()
	media.SetID("abc123")
	media.SetExtension("")

	got := media.ServeURL()
	expected := "/blog/media/abc123"
	if got != expected {
		t.Errorf("expected ServeURL %q, got %q", expected, got)
	}
}

func TestMediaHandler_DataURI(t *testing.T) {
	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 initDB(),
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	media := NewMedia().
		SetID("test123").
		SetEntityID("post1").
		SetExtension(".png").
		SetType("image/png").
		SetSize("68").
		SetURL("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==").
		SetStatus(MEDIA_STATUS_ACTIVE)

	if err := store.MediaCreate(context.Background(), media); err != nil {
		t.Fatalf("Failed to create media: %v", err)
	}

	handler := NewMediaHandler(store)

	req := httptest.NewRequest("GET", "/blog/media/test123.png", nil)
	recorder := httptest.NewRecorder()

	handler.Handler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "image/png" {
		t.Errorf("expected Content-Type image/png, got %s", recorder.Header().Get("Content-Type"))
	}

	if recorder.Header().Get("Cache-Control") != "public, max-age=31536000" {
		t.Errorf("expected Cache-Control header, got %s", recorder.Header().Get("Cache-Control"))
	}

	if recorder.Header().Get("ETag") != "test123" {
		t.Errorf("expected ETag test123, got %s", recorder.Header().Get("ETag"))
	}

	if recorder.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}

func TestMediaHandler_NotFound(t *testing.T) {
	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 initDB(),
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := NewMediaHandler(store)

	req := httptest.NewRequest("GET", "/blog/media/nonexistent.png", nil)
	recorder := httptest.NewRecorder()

	handler.Handler(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}
}

func TestMediaHandler_InactiveMedia(t *testing.T) {
	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 initDB(),
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	media := NewMedia().
		SetID("inactive123").
		SetEntityID("post1").
		SetExtension(".png").
		SetType("image/png").
		SetURL("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==").
		SetStatus(MEDIA_STATUS_DRAFT)

	if err := store.MediaCreate(context.Background(), media); err != nil {
		t.Fatalf("Failed to create media: %v", err)
	}

	handler := NewMediaHandler(store)

	req := httptest.NewRequest("GET", "/blog/media/inactive123.png", nil)
	recorder := httptest.NewRecorder()

	handler.Handler(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status %d for inactive media, got %d", http.StatusNotFound, recorder.Code)
	}
}

func TestMediaHandler_HTTPRedirect(t *testing.T) {
	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		MediaTableName:     "blog_media",
		DB:                 initDB(),
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	media := NewMedia().
		SetID("redirect123").
		SetEntityID("post1").
		SetExtension(".png").
		SetType("image/png").
		SetURL("https://example.com/image.png").
		SetStatus(MEDIA_STATUS_ACTIVE)

	if err := store.MediaCreate(context.Background(), media); err != nil {
		t.Fatalf("Failed to create media: %v", err)
	}

	handler := NewMediaHandler(store)

	req := httptest.NewRequest("GET", "/blog/media/redirect123.png", nil)
	recorder := httptest.NewRecorder()

	handler.Handler(recorder, req)

	if recorder.Code != http.StatusFound {
		t.Errorf("expected status %d for HTTP redirect, got %d", http.StatusFound, recorder.Code)
	}

	if recorder.Header().Get("Location") != "https://example.com/image.png" {
		t.Errorf("expected redirect to https://example.com/image.png, got %s", recorder.Header().Get("Location"))
	}
}
