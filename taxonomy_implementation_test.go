package blogstore

import (
	"context"
	"testing"
)

// ============================ TAXONOMY TESTS ============================

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

// ============================ TERM TESTS ============================

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

// ============================ TERM RELATION TESTS ============================

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

// ============================ STORE INTEGRATION TESTS ============================

func TestStoreTaxonomyCreateAndFind(t *testing.T) {
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

	// Create a taxonomy
	taxonomy := NewTaxonomy().
		SetName("Categories").
		SetSlug("category").
		SetDescription("Post categories")

	if err := store.TaxonomyCreate(ctx, taxonomy); err != nil {
		t.Fatalf("TaxonomyCreate() error = %v, want nil", err)
	}

	// Find by ID
	found, err := store.TaxonomyFindByID(ctx, taxonomy.GetID())
	if err != nil {
		t.Fatalf("TaxonomyFindByID() error = %v, want nil", err)
	}
	if found == nil {
		t.Fatal("TaxonomyFindByID() returned nil, want non-nil")
	}
	if found.GetName() != "Categories" {
		t.Errorf("GetName() = %q, want %q", found.GetName(), "Categories")
	}

	// Find by slug
	foundBySlug, err := store.TaxonomyFindBySlug(ctx, "category")
	if err != nil {
		t.Fatalf("TaxonomyFindBySlug() error = %v, want nil", err)
	}
	if foundBySlug == nil {
		t.Fatal("TaxonomyFindBySlug() returned nil, want non-nil")
	}
	if foundBySlug.GetName() != "Categories" {
		t.Errorf("GetName() = %q, want %q", foundBySlug.GetName(), "Categories")
	}
}

func TestStoreTaxonomyListAndCount(t *testing.T) {
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

	// Create multiple taxonomies
	taxonomies := []TaxonomyInterface{
		NewTaxonomy().SetName("Categories").SetSlug("category"),
		NewTaxonomy().SetName("Tags").SetSlug("tag"),
		NewTaxonomy().SetName("Series").SetSlug("series"),
	}

	for _, tax := range taxonomies {
		if err := store.TaxonomyCreate(ctx, tax); err != nil {
			t.Fatalf("TaxonomyCreate() error = %v, want nil", err)
		}
	}

	// Count
	count, err := store.TaxonomyCount(ctx, TaxonomyQueryOptions{})
	if err != nil {
		t.Fatalf("TaxonomyCount() error = %v, want nil", err)
	}
	if count != 3 {
		t.Errorf("TaxonomyCount() = %d, want %d", count, 3)
	}

	// List
	list, err := store.TaxonomyList(ctx, TaxonomyQueryOptions{})
	if err != nil {
		t.Fatalf("TaxonomyList() error = %v, want nil", err)
	}
	if len(list) != 3 {
		t.Errorf("TaxonomyList() len = %d, want %d", len(list), 3)
	}
}

func TestStoreTermCreateAndFind(t *testing.T) {
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

	// Create taxonomy first
	taxonomy := NewTaxonomy().SetName("Categories").SetSlug("category")
	if err := store.TaxonomyCreate(ctx, taxonomy); err != nil {
		t.Fatalf("TaxonomyCreate() error = %v, want nil", err)
	}

	// Create a term
	term := NewTerm().
		SetTaxonomyID(taxonomy.GetID()).
		SetName("Technology").
		SetSlug("technology").
		SetDescription("Tech posts")

	if err := store.TermCreate(ctx, term); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}

	// Find by ID
	found, err := store.TermFindByID(ctx, term.GetID())
	if err != nil {
		t.Fatalf("TermFindByID() error = %v, want nil", err)
	}
	if found == nil {
		t.Fatal("TermFindByID() returned nil, want non-nil")
	}
	if found.GetName() != "Technology" {
		t.Errorf("GetName() = %q, want %q", found.GetName(), "Technology")
	}

	// Find by slug
	foundBySlug, err := store.TermFindBySlug(ctx, "category", "technology")
	if err != nil {
		t.Fatalf("TermFindBySlug() error = %v, want nil", err)
	}
	if foundBySlug == nil {
		t.Fatal("TermFindBySlug() returned nil, want non-nil")
	}
	if foundBySlug.GetName() != "Technology" {
		t.Errorf("GetName() = %q, want %q", foundBySlug.GetName(), "Technology")
	}
}

func TestStoreTermListAndCount(t *testing.T) {
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

	// Create taxonomy
	taxonomy := NewTaxonomy().SetName("Categories").SetSlug("category")
	if err := store.TaxonomyCreate(ctx, taxonomy); err != nil {
		t.Fatalf("TaxonomyCreate() error = %v, want nil", err)
	}

	// Create terms
	terms := []TermInterface{
		NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("Technology").SetSlug("tech"),
		NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("Science").SetSlug("science"),
		NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("Art").SetSlug("art"),
	}

	for _, term := range terms {
		if err := store.TermCreate(ctx, term); err != nil {
			t.Fatalf("TermCreate() error = %v, want nil", err)
		}
	}

	// Count
	count, err := store.TermCount(ctx, TermQueryOptions{TaxonomyID: taxonomy.GetID()})
	if err != nil {
		t.Fatalf("TermCount() error = %v, want nil", err)
	}
	if count != 3 {
		t.Errorf("TermCount() = %d, want %d", count, 3)
	}

	// List
	list, err := store.TermList(ctx, TermQueryOptions{TaxonomyID: taxonomy.GetID()})
	if err != nil {
		t.Fatalf("TermList() error = %v, want nil", err)
	}
	if len(list) != 3 {
		t.Errorf("TermList() len = %d, want %d", len(list), 3)
	}
}

func TestStorePostTerms(t *testing.T) {
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

	// Create taxonomy and terms
	taxonomy := NewTaxonomy().SetName("Tags").SetSlug("tag")
	if err := store.TaxonomyCreate(ctx, taxonomy); err != nil {
		t.Fatalf("TaxonomyCreate() error = %v, want nil", err)
	}

	term1 := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("Go").SetSlug("go")
	term2 := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("Programming").SetSlug("programming")

	if err := store.TermCreate(ctx, term1); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}
	if err := store.TermCreate(ctx, term2); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}

	// Create a post
	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}

	// Add terms to post
	if err := store.PostTermAdd(ctx, post.GetID(), term1.GetID(), 0); err != nil {
		t.Fatalf("PostTermAdd() error = %v, want nil", err)
	}
	if err := store.PostTermAdd(ctx, post.GetID(), term2.GetID(), 1); err != nil {
		t.Fatalf("PostTermAdd() error = %v, want nil", err)
	}

	// Get terms for post
	postTerms, err := store.PostTerms(ctx, post.GetID(), "tag")
	if err != nil {
		t.Fatalf("PostTerms() error = %v, want nil", err)
	}
	if len(postTerms) != 2 {
		t.Errorf("PostTerms() len = %d, want %d", len(postTerms), 2)
	}

	// Set terms (replaces existing)
	newTerm := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("New Tag").SetSlug("new-tag")
	if err := store.TermCreate(ctx, newTerm); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}

	if err := store.PostSetTerms(ctx, post.GetID(), "tag", []string{newTerm.GetID()}); err != nil {
		t.Fatalf("PostSetTerms() error = %v, want nil", err)
	}

	// Verify terms were replaced
	postTerms, err = store.PostTerms(ctx, post.GetID(), "tag")
	if err != nil {
		t.Fatalf("PostTerms() error = %v, want nil", err)
	}
	if len(postTerms) != 1 {
		t.Errorf("PostTerms() after SetTerms len = %d, want %d", len(postTerms), 1)
	}
	if postTerms[0].GetName() != "New Tag" {
		t.Errorf("PostTerms()[0].GetName() = %q, want %q", postTerms[0].GetName(), "New Tag")
	}
}

func TestStoreTermIncrementDecrementCount(t *testing.T) {
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

	// Create taxonomy and term
	taxonomy := NewTaxonomy().SetName("Tags").SetSlug("tag")
	if err := store.TaxonomyCreate(ctx, taxonomy); err != nil {
		t.Fatalf("TaxonomyCreate() error = %v, want nil", err)
	}

	term := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("Go").SetSlug("go").SetCount(0)
	if err := store.TermCreate(ctx, term); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}

	// Increment count
	if err := store.TermIncrementCount(ctx, term.GetID()); err != nil {
		t.Fatalf("TermIncrementCount() error = %v, want nil", err)
	}

	// Verify count increased
	found, err := store.TermFindByID(ctx, term.GetID())
	if err != nil {
		t.Fatalf("TermFindByID() error = %v, want nil", err)
	}
	if found.GetCount() != 1 {
		t.Errorf("GetCount() after increment = %d, want %d", found.GetCount(), 1)
	}

	// Decrement count
	if err := store.TermDecrementCount(ctx, term.GetID()); err != nil {
		t.Fatalf("TermDecrementCount() error = %v, want nil", err)
	}

	// Verify count decreased
	found, err = store.TermFindByID(ctx, term.GetID())
	if err != nil {
		t.Fatalf("TermFindByID() error = %v, want nil", err)
	}
	if found.GetCount() != 0 {
		t.Errorf("GetCount() after decrement = %d, want %d", found.GetCount(), 0)
	}
}

func TestStoreTaxonomyUpdate(t *testing.T) {
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

	// Create taxonomy
	taxonomy := NewTaxonomy().SetName("Old Name").SetSlug("old-slug").SetDescription("Old description")
	if err := store.TaxonomyCreate(ctx, taxonomy); err != nil {
		t.Fatalf("TaxonomyCreate() error = %v, want nil", err)
	}

	// Update
	taxonomy.SetName("New Name").SetDescription("New description")
	if err := store.TaxonomyUpdate(ctx, taxonomy); err != nil {
		t.Fatalf("TaxonomyUpdate() error = %v, want nil", err)
	}

	// Verify update
	found, err := store.TaxonomyFindByID(ctx, taxonomy.GetID())
	if err != nil {
		t.Fatalf("TaxonomyFindByID() error = %v, want nil", err)
	}
	if found.GetName() != "New Name" {
		t.Errorf("GetName() = %q, want %q", found.GetName(), "New Name")
	}
	if found.GetDescription() != "New description" {
		t.Errorf("GetDescription() = %q, want %q", found.GetDescription(), "New description")
	}
}

func TestStoreTermUpdate(t *testing.T) {
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

	// Create taxonomy and term
	taxonomy := NewTaxonomy().SetName("Categories").SetSlug("category")
	if err := store.TaxonomyCreate(ctx, taxonomy); err != nil {
		t.Fatalf("TaxonomyCreate() error = %v, want nil", err)
	}

	term := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("Old Name").SetSlug("old-name")
	if err := store.TermCreate(ctx, term); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}

	// Update
	term.SetName("New Name").SetCount(10)
	if err := store.TermUpdate(ctx, term); err != nil {
		t.Fatalf("TermUpdate() error = %v, want nil", err)
	}

	// Verify update
	found, err := store.TermFindByID(ctx, term.GetID())
	if err != nil {
		t.Fatalf("TermFindByID() error = %v, want nil", err)
	}
	if found.GetName() != "New Name" {
		t.Errorf("GetName() = %q, want %q", found.GetName(), "New Name")
	}
	if found.GetCount() != 10 {
		t.Errorf("GetCount() = %d, want %d", found.GetCount(), 10)
	}
}

func TestStoreTermDelete(t *testing.T) {
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

	// Create taxonomy and term
	taxonomy := NewTaxonomy().SetName("Categories").SetSlug("category")
	if err := store.TaxonomyCreate(ctx, taxonomy); err != nil {
		t.Fatalf("TaxonomyCreate() error = %v, want nil", err)
	}

	term := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("Technology").SetSlug("tech")
	if err := store.TermCreate(ctx, term); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}

	// Delete
	if err := store.TermDelete(ctx, term); err != nil {
		t.Fatalf("TermDelete() error = %v, want nil", err)
	}

	// Verify deleted
	found, err := store.TermFindByID(ctx, term.GetID())
	if err != nil {
		t.Fatalf("TermFindByID() error = %v, want nil", err)
	}
	if found != nil {
		t.Errorf("TermFindByID() = %#v, want nil", found)
	}
}
