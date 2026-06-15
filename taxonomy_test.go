package blogstore

import (
	"testing"
)

// ============================ TAXONOMY IMPLEMENTATION TESTS ============================

func TestNewTaxonomyDefaults(t *testing.T) {
	taxonomy := NewTaxonomy()

	if taxonomy == nil {
		t.Fatalf("NewTaxonomy() returned nil")
	}

	if taxonomy.GetID() == "" {
		t.Errorf("NewTaxonomy() must set a non-empty ID")
	}

	if taxonomy.GetName() != "" {
		t.Errorf("NewTaxonomy() name = %q, want empty string", taxonomy.GetName())
	}

	if taxonomy.GetSlug() != "" {
		t.Errorf("NewTaxonomy() slug = %q, want empty string", taxonomy.GetSlug())
	}

	if taxonomy.GetCreatedAt() == "" {
		t.Errorf("NewTaxonomy() CreatedAt must not be empty")
	}

	if taxonomy.GetUpdatedAt() == "" {
		t.Errorf("NewTaxonomy() UpdatedAt must not be empty")
	}
}

func TestTaxonomySetters(t *testing.T) {
	taxonomy := NewTaxonomy()

	taxonomy.SetName("Categories").
		SetSlug("category").
		SetDescription("Post categories")

	if taxonomy.GetName() != "Categories" {
		t.Errorf("GetName() = %q, want %q", taxonomy.GetName(), "Categories")
	}

	if taxonomy.GetSlug() != "category" {
		t.Errorf("GetSlug() = %q, want %q", taxonomy.GetSlug(), "category")
	}

	if taxonomy.GetDescription() != "Post categories" {
		t.Errorf("GetDescription() = %q, want %q", taxonomy.GetDescription(), "Post categories")
	}
}

func TestTaxonomySlugNormalization(t *testing.T) {
	taxonomy := NewTaxonomy()
	taxonomy.SetSlug("My Taxonomy!")

	// Slug should be normalized
	if taxonomy.GetSlug() != "my-taxonomy" {
		t.Errorf("GetSlug() = %q, want %q", taxonomy.GetSlug(), "my-taxonomy")
	}
}

// ============================ TERM IMPLEMENTATION TESTS ============================

func TestNewTermDefaults(t *testing.T) {
	term := NewTerm()

	if term == nil {
		t.Fatalf("NewTerm() returned nil")
	}

	if term.GetID() == "" {
		t.Errorf("NewTerm() must set a non-empty ID")
	}

	if term.GetTaxonomyID() != "" {
		t.Errorf("NewTerm() taxonomy_id = %q, want empty string", term.GetTaxonomyID())
	}

	if term.GetParentID() != "" {
		t.Errorf("NewTerm() parent_id = %q, want empty string", term.GetParentID())
	}

	if term.GetSequence() != 0 {
		t.Errorf("NewTerm() sequence = %d, want 0", term.GetSequence())
	}

	if term.GetName() != "" {
		t.Errorf("NewTerm() name = %q, want empty string", term.GetName())
	}

	if term.GetSlug() != "" {
		t.Errorf("NewTerm() slug = %q, want empty string", term.GetSlug())
	}

	if term.GetCount() != 0 {
		t.Errorf("NewTerm() count = %d, want 0", term.GetCount())
	}

	if term.GetCreatedAt() == "" {
		t.Errorf("NewTerm() CreatedAt must not be empty")
	}

	if term.GetUpdatedAt() == "" {
		t.Errorf("NewTerm() UpdatedAt must not be empty")
	}
}

func TestTermSetters(t *testing.T) {
	term := NewTerm()

	term.SetTaxonomyID("taxonomy-123").
		SetParentID("parent-456").
		SetSequence(5).
		SetName("Technology").
		SetSlug("tech").
		SetDescription("Tech posts").
		SetCount(5)

	if term.GetTaxonomyID() != "taxonomy-123" {
		t.Errorf("GetTaxonomyID() = %q, want %q", term.GetTaxonomyID(), "taxonomy-123")
	}

	if term.GetParentID() != "parent-456" {
		t.Errorf("GetParentID() = %q, want %q", term.GetParentID(), "parent-456")
	}

	if term.GetSequence() != 5 {
		t.Errorf("GetSequence() = %d, want %d", term.GetSequence(), 5)
	}

	if term.GetName() != "Technology" {
		t.Errorf("GetName() = %q, want %q", term.GetName(), "Technology")
	}

	if term.GetSlug() != "tech" {
		t.Errorf("GetSlug() = %q, want %q", term.GetSlug(), "tech")
	}

	if term.GetDescription() != "Tech posts" {
		t.Errorf("GetDescription() = %q, want %q", term.GetDescription(), "Tech posts")
	}

	if term.GetCount() != 5 {
		t.Errorf("GetCount() = %d, want %d", term.GetCount(), 5)
	}
}

func TestTermSlugNormalization(t *testing.T) {
	term := NewTerm()
	term.SetSlug("My Term Name!")

	// Slug should be normalized
	if term.GetSlug() != "my-term-name" {
		t.Errorf("GetSlug() = %q, want %q", term.GetSlug(), "my-term-name")
	}
}

func TestTermSequenceManipulation(t *testing.T) {
	term := NewTerm()

	// Test setting sequence
	term.SetSequence(10)
	if term.GetSequence() != 10 {
		t.Errorf("GetSequence() = %d, want %d", term.GetSequence(), 10)
	}

	// Test setting to 0
	term.SetSequence(0)
	if term.GetSequence() != 0 {
		t.Errorf("GetSequence() = %d, want %d", term.GetSequence(), 0)
	}

	// Test negative sequence
	term.SetSequence(-5)
	if term.GetSequence() != -5 {
		t.Errorf("GetSequence() = %d, want %d", term.GetSequence(), -5)
	}
}

func TestTermCountManipulation(t *testing.T) {
	term := NewTerm()

	// Test incrementing
	term.SetCount(10)
	if term.GetCount() != 10 {
		t.Errorf("GetCount() = %d, want %d", term.GetCount(), 10)
	}

	// Test setting to 0
	term.SetCount(0)
	if term.GetCount() != 0 {
		t.Errorf("GetCount() = %d, want %d", term.GetCount(), 0)
	}
}

// ============================ TERM RELATION IMPLEMENTATION TESTS ============================

func TestNewTermRelationDefaults(t *testing.T) {
	rel := NewTermRelation()

	if rel == nil {
		t.Fatalf("NewTermRelation() returned nil")
	}

	if rel.GetID() == "" {
		t.Errorf("NewTermRelation() must set a non-empty ID")
	}

	if rel.GetPostID() != "" {
		t.Errorf("NewTermRelation() post_id = %q, want empty string", rel.GetPostID())
	}

	if rel.GetTermID() != "" {
		t.Errorf("NewTermRelation() term_id = %q, want empty string", rel.GetTermID())
	}

	if rel.GetSequence() != 0 {
		t.Errorf("NewTermRelation() sequence = %d, want 0", rel.GetSequence())
	}

	if rel.GetCreatedAt() == "" {
		t.Errorf("NewTermRelation() CreatedAt must not be empty")
	}

	if rel.GetUpdatedAt() == "" {
		t.Errorf("NewTermRelation() UpdatedAt must not be empty")
	}
}

func TestTermRelationSetters(t *testing.T) {
	rel := NewTermRelation()

	rel.SetPostID("post-123").
		SetTermID("term-456").
		SetSequence(3)

	if rel.GetPostID() != "post-123" {
		t.Errorf("GetPostID() = %q, want %q", rel.GetPostID(), "post-123")
	}

	if rel.GetTermID() != "term-456" {
		t.Errorf("GetTermID() = %q, want %q", rel.GetTermID(), "term-456")
	}

	if rel.GetSequence() != 3 {
		t.Errorf("GetSequence() = %d, want %d", rel.GetSequence(), 3)
	}
}
