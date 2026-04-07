# Taxonomy System Proposal for Blog Store

## Executive Summary

This proposal outlines a **unified taxonomy system** to add tags, categories, and future classification types to Blog Store. This architecture provides maximum flexibility while maintaining database integrity and following existing codebase patterns.

## Problem Statement

Currently, Blog Store lacks any content classification mechanism. As content grows, users need:
- **Tags**: Non-hierarchical, flat labels (e.g., "go", "tutorial", "api")
- **Categories**: Hierarchical, structured organization (e.g., "Technology > Programming > Go")
- **Future types**: Series, difficulty levels, languages, etc.

## Proposed Solution

### Three-Table Taxonomy Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  blog_taxonomy  │    │   blog_term     │    │  blog_term_rel  │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ id (PK)         │───▶│ id (PK)         │    │ id (PK)         │
│ name            │    │ taxonomy_id     │───▶│ post_id (FK)    │
│ slug (unique)   │    │ parent_id (FK)  │──┐ │ term_id (FK)    │
│ description     │    │ name            │  │ │ order           │
│ created_at      │    │ slug (unique)   │  │ │ created_at      │
└─────────────────┘    │ description     │  │ └─────────────────┘
                         │ count           │  │
                         │ created_at      │  └─ Self-referencing
                         └─────────────────┘     for hierarchy
```

### Why This Architecture?

| Feature | Benefit |
|---------|---------|
| **Unified schema** | One set of code handles all classification types |
| **Extensible** | New taxonomy = insert row, no schema changes |
| **Hierarchical** | Built-in parent/child support for categories |
| **Performant** | Denormalized `count` field avoids expensive COUNT queries |
| **Ordered relations** | `order` field enables manual term sequencing |

## Implementation Plan

### Phase 1: Core Domain (Day 1)

#### Files to Create

1. **taxonomy_interface.go** - Interface definitions
```go
type TaxonomyInterface interface {
    GetID() string
    GetName() string        // Display name (e.g., "Categories")
    GetSlug() string        // URL slug (e.g., "category")
    GetDescription() string
    SetID(id string) TaxonomyInterface
    SetName(name string) TaxonomyInterface
    SetSlug(slug string) TaxonomyInterface
}

type TermInterface interface {
    GetID() string
    GetTaxonomyID() string
    GetParentID() string    // For hierarchy (empty if root)
    GetName() string
    GetSlug() string
    GetDescription() string
    GetCount() int          // Cached post count
    // ... setters, Carbon timestamps, DataObject methods
}

type TermRelationInterface interface {
    GetID() string
    GetPostID() string
    GetTermID() string
    GetSequence() int       // For manual ordering (0 = default)
    // ... setters, timestamps
}
```

2. **taxonomy_implementation.go** - Concrete types using `dataobject.DataObject`

3. **taxonomy_query_options.go** - Query options for filtering terms
```go
type TermQueryOptions struct {
    ID           string
    TaxonomyID   string
    TaxonomySlug string
    ParentID     string
    Search       string
    Limit        int
    Offset       int
    OrderBy      string
    SortOrder    string
}

type TermRelationQueryOptions struct {
    PostID       string
    TermID       string
    TaxonomyID   string
    TaxonomySlug string
}
```

### Phase 2: Store Layer (Day 2-3)

#### Extend StoreInterface
```go
type StoreInterface interface {
    // ... existing methods ...
    
    // Taxonomy management
    TaxonomyCount(ctx context.Context, options TaxonomyQueryOptions) (int64, error)
    TaxonomyCreate(ctx context.Context, taxonomy TaxonomyInterface) error
    TaxonomyDelete(ctx context.Context, taxonomy TaxonomyInterface) error
    TaxonomyFindByID(ctx context.Context, id string) (TaxonomyInterface, error)
    TaxonomyFindBySlug(ctx context.Context, slug string) (TaxonomyInterface, error)
    TaxonomyList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyInterface, error)
    TaxonomyUpdate(ctx context.Context, taxonomy TaxonomyInterface) error
    
    // Term management
    TermCount(ctx context.Context, options TermQueryOptions) (int64, error)
    TermCreate(ctx context.Context, term TermInterface) error
    TermDelete(ctx context.Context, term TermInterface) error
    TermFindByID(ctx context.Context, id string) (TermInterface, error)
    TermFindBySlug(ctx context.Context, taxonomySlug, termSlug string) (TermInterface, error)
    TermList(ctx context.Context, options TermQueryOptions) ([]TermInterface, error)
    TermUpdate(ctx context.Context, term TermInterface) error
    
    // Post-term relationships
    PostTermAdd(ctx context.Context, postID string, termID string, order int) error
    PostTermRemove(ctx context.Context, postID string, termID string) error
    PostTerms(ctx context.Context, postID string, taxonomySlug string) ([]TermInterface, error)
    PostSetTerms(ctx context.Context, postID string, taxonomySlug string, termIDs []string) error
    
    // Utility queries
    TermIncrementCount(ctx context.Context, termID string) error
    TermDecrementCount(ctx context.Context, termID string) error
}
```

#### Implementation in taxonomy_store.go

**Key design decisions:**

1. **Slug uniqueness**: `(taxonomy_id, slug)` composite unique index on `blog_term`
   - Allows "programming" in both "category" and "tag" taxonomies
   - Prevents duplicates within same taxonomy

2. **Count caching**: `count` field updated via triggers or application logic
   - Avoids expensive `COUNT(*)` on post listings
   - Recalculated when posts are added/removed

3. **Self-referencing hierarchy**: `parent_id` references `blog_term.id`
   - Enables unlimited nesting depth
   - Root terms have `parent_id = ""`

### Phase 3: Database Schema (Day 2)

#### Add to sql_create_table.go

```go
func (store *storeImplementation) sqlCreateTaxonomyTable() (string, error) {
    // blog_taxonomy: category, tag, series, etc.
}

func (store *storeImplementation) sqlCreateTermTable() (string, error) {
    // blog_term: actual terms with hierarchy
}

func (store *storeImplementation) sqlCreateTermRelationTable() (string, error) {
    // blog_term_rel: post-to-term associations
}
```

**SQLite DDL:**

```sql
-- Taxonomy types (category, tag, etc.)
CREATE TABLE IF NOT EXISTS blog_taxonomy (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Terms within taxonomies
CREATE TABLE IF NOT EXISTS blog_term (
    id VARCHAR(32) PRIMARY KEY,
    taxonomy_id VARCHAR(32) NOT NULL,
    parent_id VARCHAR(32) DEFAULT '',
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL,
    description TEXT,
    count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(taxonomy_id, slug),
    FOREIGN KEY (taxonomy_id) REFERENCES blog_taxonomy(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES blog_term(id) ON DELETE SET NULL
);

-- Post-term relationships
CREATE TABLE IF NOT EXISTS blog_term_rel (
    id VARCHAR(32) PRIMARY KEY,
    post_id VARCHAR(32) NOT NULL,
    term_id VARCHAR(32) NOT NULL,
    sequence INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(post_id, term_id),
    FOREIGN KEY (post_id) REFERENCES blog_post(id) ON DELETE CASCADE,
    FOREIGN KEY (term_id) REFERENCES blog_term(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_term_taxonomy ON blog_term(taxonomy_id);
CREATE INDEX IF NOT EXISTS idx_term_parent ON blog_term(parent_id);
CREATE INDEX IF NOT EXISTS idx_term_rel_post ON blog_term_rel(post_id);
CREATE INDEX IF NOT EXISTS idx_term_rel_term ON blog_term_rel(term_id);
```

### Phase 4: Post Integration (Day 3)

#### Extend PostInterface

```go
type PostInterface interface {
    // ... existing methods ...
    
    // Taxonomy methods
    TermIDs(taxonomySlug string) []string
    SetTermIDs(taxonomySlug string, termIDs []string) PostInterface
    
    // Convenience helpers
    CategoryIDs() []string
    SetCategoryIDs(ids []string) PostInterface
    TagIDs() []string
    SetTagIDs(ids []string) PostInterface
}
```

#### Query Options Extension

```go
type PostQueryOptions struct {
    // ... existing fields ...
    
    TermFilter       map[string]string  // taxonomy_slug => term_slug
    TermIDFilter     []string           // term IDs
    WithTerms        bool               // include terms in response
}
```

### Phase 5: MCP Tools (Day 4)

#### New Tools

```json
{
  "name": "taxonomy_list",
  "description": "List available taxonomy types (category, tag, etc.)"
},
{
  "name": "taxonomy_create",
  "description": "Create a new taxonomy type",
  "inputSchema": {
    "required": ["name", "slug"],
    "properties": {
      "name": {"type": "string", "description": "Display name"},
      "slug": {"type": "string", "description": "URL slug"},
      "description": {"type": "string"}
    }
  }
},
{
  "name": "term_list",
  "description": "List terms within a taxonomy",
  "inputSchema": {
    "properties": {
      "taxonomy": {"type": "string", "description": "Taxonomy slug (category, tag)"},
      "parent_id": {"type": "string", "description": "Filter by parent term"},
      "search": {"type": "string"},
      "limit": {"type": "integer"},
      "offset": {"type": "integer"}
    }
  }
},
{
  "name": "term_create",
  "description": "Create a term",
  "inputSchema": {
    "required": ["taxonomy", "name"],
    "properties": {
      "taxonomy": {"type": "string"},
      "name": {"type": "string"},
      "slug": {"type": "string"},
      "parent_id": {"type": "string", "description": "For hierarchical taxonomies"},
      "description": {"type": "string"}
    }
  }
},
{
  "name": "post_set_terms",
  "description": "Set terms for a post (replaces existing)",
  "inputSchema": {
    "required": ["post_id", "taxonomy"],
    "properties": {
      "post_id": {"type": "string"},
      "taxonomy": {"type": "string"},
      "terms": {"type": "array", "items": {"type": "string"}, "description": "Term slugs or IDs"}
    }
  }
},
{
  "name": "post_add_term",
  "description": "Add a single term to a post",
  "inputSchema": {
    "required": ["post_id", "taxonomy", "term"],
    "properties": {
      "post_id": {"type": "string"},
      "taxonomy": {"type": "string"},
      "term": {"type": "string"}
    }
  }
}
```

### Phase 6: Constants & Utilities (Day 1)

#### Add to constants.go

```go
// Taxonomy types (pre-defined)
const TAXONOMY_CATEGORY = "category"
const TAXONOMY_TAG = "tag"

// Table columns
const COLUMN_TAXONOMY_ID = "taxonomy_id"
const COLUMN_PARENT_ID = "parent_id"
const COLUMN_TERM_ID = "term_id"
const COLUMN_TERM_SEQUENCE = "sequence"
```

## API Usage Examples

### Creating Default Taxonomies

```go
// Store initialization creates default taxonomies
func (store *storeImplementation) createDefaultTaxonomies() error {
    // Create "category" taxonomy
    cat := NewTaxonomy()
    cat.SetName("Categories").SetSlug(TAXONOMY_CATEGORY)
    store.TaxonomyCreate(ctx, cat)
    
    // Create "tag" taxonomy
    tag := NewTaxonomy()
    tag.SetName("Tags").SetSlug(TAXONOMY_TAG)
    store.TaxonomyCreate(ctx, tag)
}
```

### Working with Categories

```go
// Create hierarchical categories
tech := NewTerm()
tech.SetTaxonomyID("category").
    SetName("Technology").
    SetSlug("technology")
store.TermCreate(ctx, tech)

prog := NewTerm()
prog.SetTaxonomyID("category").
    SetParentID(tech.GetID()).
    SetName("Programming").
    SetSlug("programming")
store.TermCreate(ctx, prog)

// Assign to post
store.PostSetTerms(ctx, post.GetID(), TAXONOMY_CATEGORY, []string{prog.GetID()})
```

### Working with Tags

```go
// Create tags
goTag := NewTerm()
goTag.SetTaxonomyID("tag").SetName("Go").SetSlug("go")
store.TermCreate(ctx, goTag)

// Assign multiple tags
store.PostSetTerms(ctx, post.GetID(), TAXONOMY_TAG, []string{
    goTag.GetID(),
    "programming",
    "tutorial",
})
```

### Querying Posts by Term

```go
// Get all posts in "programming" category
posts, _ := store.PostList(ctx, PostQueryOptions{
    TermFilter: map[string]string{
        "taxonomy": "category",
        "slug": "programming",
    },
})

// Get posts with any of these tags
posts, _ := store.PostList(ctx, PostQueryOptions{
    TermIDFilter: []string{"tag1", "tag2", "tag3"},
})
```

## Migration Strategy

### For Existing Installations

1. **Safe migration**: Taxonomy tables are additive, no existing data changes
2. **Optional feature**: Taxonomy methods return empty if tables don't exist
3. **Backward compatible**: All existing code continues to work

### Migration Helper

```go
func (store *storeImplementation) migrateTaxonomies() error {
    // Create tables
    // Insert default taxonomies (category, tag)
    // No changes to existing posts
}
```

## Testing Strategy

1. **Unit tests**: `taxonomy_test.go` - CRUD operations
2. **Integration tests**: `taxonomy_store_test.go` - Store methods
3. **Hierarchy tests**: Parent/child relationships, cycle prevention
4. **MCP tests**: Tool validation in `mcp_test.go`

## Timeline

| Phase | Duration | Deliverable |
|-------|----------|-------------|
| 1. Domain models | 1 day | taxonomy_interface.go, taxonomy_implementation.go |
| 2. Store layer | 2 days | taxonomy_store.go, StoreInterface extensions |
| 3. Database schema | 1 day | SQL DDL, automigration |
| 4. Post integration | 1 day | Post terms methods |
| 5. MCP tools | 1 day | 6 new MCP tools |
| 6. Tests & polish | 1 day | Full test coverage |

**Total: 7 days**

## Future Extensibility

This architecture supports:
- **Custom taxonomies**: "series", "difficulty", "language"
- **Term metadata**: JSON field for icons, colors, descriptions
- **Term templates**: Default content for new posts in term
- **Term permissions**: Who can create terms
- **Related terms**: Cross-taxonomy associations
- **Term subscriptions**: Users following terms

## Recommendation

**Approve the taxonomy system approach.** It provides:
- ✅ Single unified architecture for all classification needs
- ✅ No schema changes when adding new taxonomy types
- ✅ Built-in hierarchy support
- ✅ Follows existing Blog Store patterns
- ✅ Full MCP integration
- ✅ Maintains referential integrity
