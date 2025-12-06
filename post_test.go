package blogstore

import (
	"encoding/json"
	"testing"

	"github.com/dracory/sb"
)

// TestNewPostDefaults tests that NewPost() returns a Post with:
// - a non-empty ID,
// - status set to DRAFT
// - non-empty PublishedAt,
// - CreatedAt, UpdatedAt, SoftDeletedAt set to MAX_DATETIME
// - Featured set to NO,
// - Metas set to empty map.
func TestNewPostDefaults(t *testing.T) {
	p := NewPost()

	if p == nil {
		t.Fatalf("NewPost() returned nil")
	}

	if p.ID() == "" {
		t.Errorf("NewPost() must set a non-empty ID")
	}

	if got := p.Status(); got != POST_STATUS_DRAFT {
		t.Errorf("NewPost() status = %q, want %q", got, POST_STATUS_DRAFT)
	}

	if got := p.PublishedAt(); got == "" {
		t.Errorf("NewPost() PublishedAt must not be empty")
	}

	if got := p.CreatedAt(); got == "" {
		t.Errorf("NewPost() CreatedAt must not be empty")
	}

	if got := p.UpdatedAt(); got == "" {
		t.Errorf("NewPost() UpdatedAt must not be empty")
	}

	if got := p.SoftDeletedAt(); got != sb.MAX_DATETIME {
		t.Errorf("NewPost() SoftDeletedAt = %q, want %q", got, sb.MAX_DATETIME)
	}

	if got := p.Featured(); got != NO {
		t.Errorf("NewPost() Featured = %q, want %q", got, NO)
	}

	metas, err := p.Metas()
	if err != nil {
		t.Fatalf("NewPost() Metas error = %v, want nil", err)
	}
	if metas == nil {
		t.Fatalf("NewPost() Metas must not be nil")
	}
	if len(metas) != 0 {
		t.Errorf("NewPost() Metas length = %d, want 0", len(metas))
	}
}

func TestPostMetasAndMetaHelpers(t *testing.T) {
	p := NewPost()

	// SetMetas and Metas
	m := map[string]string{"k1": "v1"}
	if err := p.SetMetas(m); err != nil {
		t.Fatalf("SetMetas() error = %v, want nil", err)
	}

	metas, err := p.Metas()
	if err != nil {
		t.Fatalf("Metas() error = %v, want nil", err)
	}

	if got := metas["k1"]; got != "v1" {
		t.Errorf("Metas()[k1] = %q, want %q", got, "v1")
	}

	// AddMetas merges keys
	add := map[string]string{"k2": "v2"}
	if err := p.AddMetas(add); err != nil {
		t.Fatalf("AddMetas() error = %v, want nil", err)
	}

	metas, err = p.Metas()
	if err != nil {
		t.Fatalf("Metas() after AddMetas error = %v, want nil", err)
	}

	if got := metas["k1"]; got != "v1" {
		t.Errorf("after AddMetas, metas[k1] = %q, want %q", got, "v1")
	}
	if got := metas["k2"]; got != "v2" {
		t.Errorf("after AddMetas, metas[k2] = %q, want %q", got, "v2")
	}

	// Meta / SetMeta helpers
	if err := p.SetMeta("editor", "alice"); err != nil {
		t.Fatalf("SetMeta() error = %v, want nil", err)
	}

	if got := p.Meta("editor"); got != "alice" {
		t.Errorf("Meta(\"editor\") = %q, want %q", got, "alice")
	}

	if got := p.Editor(); got != "alice" {
		t.Errorf("Editor() = %q, want %q", got, "alice")
	}
}

func TestPostMetasJSONRoundTrip(t *testing.T) {
	p := NewPost()

	m := map[string]string{"foo": "bar"}
	if err := p.SetMetas(m); err != nil {
		t.Fatalf("SetMetas() error = %v, want nil", err)
	}

	raw := p.Get(COLUMN_METAS)
	if raw == "" {
		t.Fatalf("COLUMN_METAS must not be empty after SetMetas")
	}

	decoded := map[string]string{}
	if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
		t.Fatalf("stored metas is not valid JSON: %v", err)
	}

	if got := decoded["foo"]; got != "bar" {
		t.Errorf("decoded[foo] = %q, want %q", got, "bar")
	}
}

func TestPostStatusHelpers(t *testing.T) {
	p := NewPost()

	// default is draft
	if !p.IsDraft() {
		t.Errorf("IsDraft() = false, want true for default post")
	}
	if p.IsPublished() {
		t.Errorf("IsPublished() = true, want false for default post")
	}
	if p.IsTrashed() {
		t.Errorf("IsTrashed() = true, want false for default post")
	}

	// published
	p.SetStatus(POST_STATUS_PUBLISHED)
	if !p.IsPublished() {
		t.Errorf("IsPublished() = false, want true for published status")
	}
	if p.IsDraft() {
		t.Errorf("IsDraft() = true, want false for published status")
	}
	if p.IsUnpublished() {
		t.Errorf("IsUnpublished() = true, want false for published status")
	}

	// trashed
	p.SetStatus(POST_STATUS_TRASH)
	if !p.IsTrashed() {
		t.Errorf("IsTrashed() = false, want true for trash status")
	}
}

func TestPostSlugAndImageUrlOrDefault(t *testing.T) {
	p := NewPost()

	p.SetTitle("Hello World Post")
	if got := p.Slug(); got == "" {
		t.Errorf("Slug() must not be empty")
	}

	// default image uses fallback URL
	if got := p.ImageUrl(); got != "" {
		t.Errorf("default ImageUrl() = %q, want empty", got)
	}
	if got := p.ImageUrlOrDefault(); got == "" {
		t.Errorf("ImageUrlOrDefault() must not be empty")
	}

	// when image is set, ImageUrlOrDefault returns provided value
	p.SetImageUrl("http://example.com/img.png")
	if got := p.ImageUrlOrDefault(); got != "http://example.com/img.png" {
		t.Errorf("ImageUrlOrDefault() = %q, want %q", got, "http://example.com/img.png")
	}
}
