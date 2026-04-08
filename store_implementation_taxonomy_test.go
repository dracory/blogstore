package blogstore

import (
	"context"
	"testing"

	"github.com/samber/lo"
)

// ============================ TAXONOMY STORE TESTS ============================

func TestStoreTaxonomyCreateAndFind(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
		TaxonomyEnabled:    true,
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
		TaxonomyEnabled:    true,
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

func TestStoreTaxonomyUpdate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
		TaxonomyEnabled:    true,
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

// ============================ TERM STORE TESTS ============================

func TestStoreTermCreateAndFind(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
		TaxonomyEnabled:    true,
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
		TaxonomyEnabled:    true,
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

func TestStoreTermUpdate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
		TaxonomyEnabled:    true,
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
		TaxonomyEnabled:    true,
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

// ============================ POST-TERM RELATIONSHIP STORE TESTS ============================

func TestStoreTermListByPostID(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
		TaxonomyEnabled:    true,
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
	if err := store.PostTermAddAt(ctx, post.GetID(), term1.GetID(), 0); err != nil {
		t.Fatalf("PostTermAddAt() error = %v, want nil", err)
	}
	if err := store.PostTermAddAt(ctx, post.GetID(), term2.GetID(), 1); err != nil {
		t.Fatalf("PostTermAddAt() error = %v, want nil", err)
	}

	// Get terms for post
	postTerms, err := store.TermListByPostID(ctx, post.GetID(), "tag")
	if err != nil {
		t.Fatalf("TermListByPostID() error = %v, want nil", err)
	}
	if len(postTerms) != 2 {
		t.Errorf("TermListByPostID() len = %d, want %d", len(postTerms), 2)
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
	postTerms, err = store.TermListByPostID(ctx, post.GetID(), "tag")
	if err != nil {
		t.Fatalf("TermListByPostID() error = %v, want nil", err)
	}
	if len(postTerms) != 1 {
		t.Errorf("TermListByPostID() after SetTerms len = %d, want %d", len(postTerms), 1)
	}
	if postTerms[0].GetName() != "New Tag" {
		t.Errorf("TermListByPostID()[0].GetName() = %q, want %q", postTerms[0].GetName(), "New Tag")
	}
}

func TestStorePostListByTermID(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
		TaxonomyEnabled:    true,
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

	term := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("Go").SetSlug("go")
	if err := store.TermCreate(ctx, term); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}

	// Create posts
	post1 := NewPost().SetTitle("Post 1").SetStatus(POST_STATUS_PUBLISHED)
	post2 := NewPost().SetTitle("Post 2").SetStatus(POST_STATUS_PUBLISHED)
	post3 := NewPost().SetTitle("Post 3").SetStatus(POST_STATUS_PUBLISHED)

	if err := store.PostCreate(ctx, post1); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}
	if err := store.PostCreate(ctx, post2); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}
	if err := store.PostCreate(ctx, post3); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}

	// Add term to post1 and post2
	if err := store.PostTermAddAt(ctx, post1.GetID(), term.GetID(), 0); err != nil {
		t.Fatalf("PostTermAddAt() error = %v, want nil", err)
	}
	if err := store.PostTermAddAt(ctx, post2.GetID(), term.GetID(), 1); err != nil {
		t.Fatalf("PostTermAddAt() error = %v, want nil", err)
	}

	// List posts for this term
	posts, err := store.PostListByTermID(ctx, term.GetID(), PostQueryOptions{})
	if err != nil {
		t.Fatalf("PostListByTermID() error = %v, want nil", err)
	}
	if len(posts) != 2 {
		t.Errorf("PostListByTermID() len = %d, want %d", len(posts), 2)
	}

	// Verify correct posts returned
	postIDs := []string{posts[0].GetID(), posts[1].GetID()}
	if !lo.Contains(postIDs, post1.GetID()) {
		t.Errorf("PostListByTermID() should contain post1")
	}
	if !lo.Contains(postIDs, post2.GetID()) {
		t.Errorf("PostListByTermID() should contain post2")
	}
	if lo.Contains(postIDs, post3.GetID()) {
		t.Errorf("PostListByTermID() should NOT contain post3")
	}
}

func TestStoreTermIncrementDecrementCount(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
		TaxonomyEnabled:    true,
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

func TestStoreTermHierarchyWithSequence(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
		TaxonomyEnabled:    true,
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

	// Create parent category
	parent := NewTerm().
		SetTaxonomyID(taxonomy.GetID()).
		SetName("Electronics").
		SetSlug("electronics")
	if err := store.TermCreate(ctx, parent); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}

	// Create subcategories with sequences
	sub1 := NewTerm().
		SetTaxonomyID(taxonomy.GetID()).
		SetParentID(parent.GetID()).
		SetSequence(3).
		SetName("Phones").
		SetSlug("phones")
	sub2 := NewTerm().
		SetTaxonomyID(taxonomy.GetID()).
		SetParentID(parent.GetID()).
		SetSequence(1).
		SetName("Laptops").
		SetSlug("laptops")
	sub3 := NewTerm().
		SetTaxonomyID(taxonomy.GetID()).
		SetParentID(parent.GetID()).
		SetSequence(2).
		SetName("Tablets").
		SetSlug("tablets")

	if err := store.TermCreate(ctx, sub1); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}
	if err := store.TermCreate(ctx, sub2); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}
	if err := store.TermCreate(ctx, sub3); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}

	// Verify subcategories have correct sequences
	found1, err := store.TermFindByID(ctx, sub1.GetID())
	if err != nil {
		t.Fatalf("TermFindByID() error = %v, want nil", err)
	}
	if found1.GetSequence() != 3 {
		t.Errorf("sub1.GetSequence() = %d, want %d", found1.GetSequence(), 3)
	}

	found2, err := store.TermFindByID(ctx, sub2.GetID())
	if err != nil {
		t.Fatalf("TermFindByID() error = %v, want nil", err)
	}
	if found2.GetSequence() != 1 {
		t.Errorf("sub2.GetSequence() = %d, want %d", found2.GetSequence(), 1)
	}

	found3, err := store.TermFindByID(ctx, sub3.GetID())
	if err != nil {
		t.Fatalf("TermFindByID() error = %v, want nil", err)
	}
	if found3.GetSequence() != 2 {
		t.Errorf("sub3.GetSequence() = %d, want %d", found3.GetSequence(), 2)
	}

	// Verify parent ID is set correctly
	if found1.GetParentID() != parent.GetID() {
		t.Errorf("sub1.GetParentID() = %q, want %q", found1.GetParentID(), parent.GetID())
	}
}
