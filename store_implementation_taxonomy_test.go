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
	if err := store.PostInsertTermAt(ctx, post.GetID(), term1.GetID(), 0); err != nil {
		t.Fatalf("PostInsertTermAt() error = %v, want nil", err)
	}
	if err := store.PostInsertTermAt(ctx, post.GetID(), term2.GetID(), 1); err != nil {
		t.Fatalf("PostInsertTermAt() error = %v, want nil", err)
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
	if err := store.PostInsertTermAt(ctx, post1.GetID(), term.GetID(), 0); err != nil {
		t.Fatalf("PostInsertTermAt() error = %v, want nil", err)
	}
	if err := store.PostInsertTermAt(ctx, post2.GetID(), term.GetID(), 1); err != nil {
		t.Fatalf("PostInsertTermAt() error = %v, want nil", err)
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

// ============================ POST-TERM SEQUENCE TESTS ============================

// Helper function to get term sequences for a post
func getPostTermSequences(t *testing.T, ctx context.Context, store StoreInterface, postID string) map[string]int {
	t.Helper()

	// Use the store's db directly to query term relations
	sqlStr := "SELECT term_id, sequence FROM " + store.(*storeImplementation).termRelationTableName + " WHERE post_id = ? ORDER BY sequence"
	rows, err := store.(*storeImplementation).db.QueryContext(ctx, sqlStr, postID)
	if err != nil {
		t.Fatalf("Failed to query term sequences: %v", err)
	}
	defer rows.Close()

	sequences := make(map[string]int)
	for rows.Next() {
		var termID string
		var seq int
		if err := rows.Scan(&termID, &seq); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		sequences[termID] = seq
	}
	return sequences
}

func TestStorePostAddTerm(t *testing.T) {
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
	term2 := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("Rust").SetSlug("rust")
	term3 := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("Python").SetSlug("python")

	if err := store.TermCreate(ctx, term1); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}
	if err := store.TermCreate(ctx, term2); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}
	if err := store.TermCreate(ctx, term3); err != nil {
		t.Fatalf("TermCreate() error = %v, want nil", err)
	}

	// Create a post
	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}

	// Add first term - should get sequence 1 (max is 0, so 0+1)
	if err := store.PostAddTerm(ctx, post.GetID(), term1.GetID()); err != nil {
		t.Fatalf("PostAddTerm() error = %v, want nil", err)
	}

	seqs := getPostTermSequences(t, ctx, store, post.GetID())
	if seqs[term1.GetID()] != 1 {
		t.Errorf("First term sequence = %d, want 1", seqs[term1.GetID()])
	}

	// Add second term - should get sequence 2
	if err := store.PostAddTerm(ctx, post.GetID(), term2.GetID()); err != nil {
		t.Fatalf("PostAddTerm() error = %v, want nil", err)
	}

	seqs = getPostTermSequences(t, ctx, store, post.GetID())
	if seqs[term2.GetID()] != 2 {
		t.Errorf("Second term sequence = %d, want 2", seqs[term2.GetID()])
	}

	// Add third term - should get sequence 3
	if err := store.PostAddTerm(ctx, post.GetID(), term3.GetID()); err != nil {
		t.Fatalf("PostAddTerm() error = %v, want nil", err)
	}

	seqs = getPostTermSequences(t, ctx, store, post.GetID())
	if seqs[term3.GetID()] != 3 {
		t.Errorf("Third term sequence = %d, want 3", seqs[term3.GetID()])
	}

	// Verify all sequences are correct
	if len(seqs) != 3 {
		t.Errorf("Total term count = %d, want 3", len(seqs))
	}

	// Verify term counts were incremented
	found, _ := store.TermFindByID(ctx, term1.GetID())
	if found.GetCount() != 1 {
		t.Errorf("term1 count = %d, want 1", found.GetCount())
	}
}

func TestStorePostAddTermErrors(t *testing.T) {
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

	// Test with empty post ID
	if err := store.PostAddTerm(ctx, "", "term-id"); err == nil {
		t.Error("PostAddTerm() with empty postID error = nil, want error")
	}

	// Test with empty term ID
	if err := store.PostAddTerm(ctx, "post-id", ""); err == nil {
		t.Error("PostAddTerm() with empty termID error = nil, want error")
	}
}

func TestStorePostMoveTermTo(t *testing.T) {
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

	termA := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("A").SetSlug("a")
	termB := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("B").SetSlug("b")
	termC := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("C").SetSlug("c")
	termD := NewTerm().SetTaxonomyID(taxonomy.GetID()).SetName("D").SetSlug("d")

	for _, term := range []TermInterface{termA, termB, termC, termD} {
		if err := store.TermCreate(ctx, term); err != nil {
			t.Fatalf("TermCreate() error = %v, want nil", err)
		}
	}

	// Create a post
	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}

	// Add terms at specific positions: A(1), B(2), C(3), D(4)
	for i, term := range []TermInterface{termA, termB, termC, termD} {
		if err := store.PostInsertTermAt(ctx, post.GetID(), term.GetID(), i+1); err != nil {
			t.Fatalf("PostInsertTermAt() error = %v, want nil", err)
		}
	}

	// Move B (pos 2) to position 4 (move down)
	// Expected: A(1), C(2), D(3), B(4)
	if err := store.PostMoveTermTo(ctx, post.GetID(), termB.GetID(), 4); err != nil {
		t.Fatalf("PostMoveTermTo() error = %v, want nil", err)
	}

	seqs := getPostTermSequences(t, ctx, store, post.GetID())
	if seqs[termA.GetID()] != 1 {
		t.Errorf("After moving B to 4: A sequence = %d, want 1", seqs[termA.GetID()])
	}
	if seqs[termC.GetID()] != 2 {
		t.Errorf("After moving B to 4: C sequence = %d, want 2", seqs[termC.GetID()])
	}
	if seqs[termD.GetID()] != 3 {
		t.Errorf("After moving B to 4: D sequence = %d, want 3", seqs[termD.GetID()])
	}
	if seqs[termB.GetID()] != 4 {
		t.Errorf("After moving B to 4: B sequence = %d, want 4", seqs[termB.GetID()])
	}

	// Move D (pos 3) to position 1 (move up)
	// Current: A(1), C(2), D(3), B(4)
	// Expected: D(1), A(2), C(3), B(4)
	if err := store.PostMoveTermTo(ctx, post.GetID(), termD.GetID(), 1); err != nil {
		t.Fatalf("PostMoveTermTo() error = %v, want nil", err)
	}

	seqs = getPostTermSequences(t, ctx, store, post.GetID())
	if seqs[termD.GetID()] != 1 {
		t.Errorf("After moving D to 1: D sequence = %d, want 1", seqs[termD.GetID()])
	}
	if seqs[termA.GetID()] != 2 {
		t.Errorf("After moving D to 1: A sequence = %d, want 2", seqs[termA.GetID()])
	}
	if seqs[termC.GetID()] != 3 {
		t.Errorf("After moving D to 1: C sequence = %d, want 3", seqs[termC.GetID()])
	}
	if seqs[termB.GetID()] != 4 {
		t.Errorf("After moving D to 1: B sequence = %d, want 4", seqs[termB.GetID()])
	}

	// Move to same position (should be no-op)
	if err := store.PostMoveTermTo(ctx, post.GetID(), termA.GetID(), 2); err != nil {
		t.Fatalf("PostMoveTermTo() same position error = %v, want nil", err)
	}

	seqs = getPostTermSequences(t, ctx, store, post.GetID())
	if seqs[termA.GetID()] != 2 {
		t.Errorf("After no-op move: A sequence = %d, want 2", seqs[termA.GetID()])
	}
}

func TestStorePostMoveTermToLargeGap(t *testing.T) {
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

	// Create taxonomy and many terms
	taxonomy := NewTaxonomy().SetName("Tags").SetSlug("tag")
	if err := store.TaxonomyCreate(ctx, taxonomy); err != nil {
		t.Fatalf("TaxonomyCreate() error = %v, want nil", err)
	}

	// Create 10 terms
	terms := make([]TermInterface, 10)
	for i := 0; i < 10; i++ {
		term := NewTerm().
			SetTaxonomyID(taxonomy.GetID()).
			SetName(string(rune('A' + i))).
			SetSlug(string(rune('a' + i)))
		if err := store.TermCreate(ctx, term); err != nil {
			t.Fatalf("TermCreate() error = %v, want nil", err)
		}
		terms[i] = term
	}

	// Create a post
	post := NewPost().SetTitle("Test Post").SetStatus(POST_STATUS_PUBLISHED)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error = %v, want nil", err)
	}

	// Add terms at positions 1-10
	for i, term := range terms {
		if err := store.PostInsertTermAt(ctx, post.GetID(), term.GetID(), i+1); err != nil {
			t.Fatalf("PostInsertTermAt() error = %v, want nil", err)
		}
	}

	// Move term at position 10 (J) to position 2
	// Expected: A(1), J(2), B(3), C(4), D(5), E(6), F(7), G(8), H(9), I(10)
	if err := store.PostMoveTermTo(ctx, post.GetID(), terms[9].GetID(), 2); err != nil {
		t.Fatalf("PostMoveTermTo() error = %v, want nil", err)
	}

	seqs := getPostTermSequences(t, ctx, store, post.GetID())

	// Verify J is at position 2
	if seqs[terms[9].GetID()] != 2 {
		t.Errorf("J sequence = %d, want 2", seqs[terms[9].GetID()])
	}

	// Verify A is still at 1
	if seqs[terms[0].GetID()] != 1 {
		t.Errorf("A sequence = %d, want 1", seqs[terms[0].GetID()])
	}

	// Verify B-I shifted from 2-9 to 3-10
	for i := 1; i <= 8; i++ {
		if seqs[terms[i].GetID()] != i+2 {
			t.Errorf("Term %d sequence = %d, want %d", i, seqs[terms[i].GetID()], i+2)
		}
	}
}

func TestStorePostMoveTermToErrors(t *testing.T) {
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

	// Test with empty post ID
	if err := store.PostMoveTermTo(ctx, "", "term-id", 1); err == nil {
		t.Error("PostMoveTermTo() with empty postID error = nil, want error")
	}

	// Test with empty term ID
	if err := store.PostMoveTermTo(ctx, "post-id", "", 1); err == nil {
		t.Error("PostMoveTermTo() with empty termID error = nil, want error")
	}

	// Test with negative sequence
	if err := store.PostMoveTermTo(ctx, "post-id", "term-id", -1); err == nil {
		t.Error("PostMoveTermTo() with negative sequence error = nil, want error")
	}

	// Test moving non-existent term
	if err := store.PostMoveTermTo(ctx, "post-id", "non-existent-term", 1); err == nil {
		t.Error("PostMoveTermTo() with non-existent term error = nil, want error")
	}
}
