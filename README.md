# Blog Store

[![Tests Status](https://github.com/dracory/blogstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/blogstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/blogstore)](https://goreportcard.com/report/github.com/dracory/blogstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/blogstore)](https://pkg.go.dev/github.com/dracory/blogstore)

Stores blog posts to a database table.

## Installation

```
go get -u github.com/dracory/blogstore
```

## Setup

```go
blogStore = blogstore.NewStore(blogstore.NewStoreOptions{
	DB:                 databaseInstance,
	PostTableName:     "blog_post",
	AutomigrateEnabled: true,
	DebugEnabled:       false,
})
```

## Usage

### Creating a Post

```go
post := blogstore.NewPost()
post.SetTitle("My First Post").
    SetContent("Post content here...").
    SetStatus(blogstore.STATUS_PUBLISHED)

err := blogStore.PostCreate(ctx, post)
```

### Listing Posts

```go
posts, err := blogStore.PostList(ctx, blogstore.PostQueryOptions{
    Status: blogstore.STATUS_PUBLISHED,
    Limit:  10,
    Offset: 0,
})
```

## MCP (Model Context Protocol)

Blog Store includes an MCP (Model Context Protocol) HTTP handler that allows LLM clients (for example Windsurf) to manage blog posts via JSON-RPC tools.

- The MCP handler lives in the `mcp` package
- It supports MCP JSON-RPC methods (`initialize`, `tools/list`, `tools/call`) and legacy aliases (`list_tools`, `call_tool`)
- It exposes tools such as `post_list`, `post_create`, `post_get`, `post_update`, and `post_delete`

See the detailed documentation and examples in: `mcp/README.md`

## Taxonomy System

Blog Store includes a flexible taxonomy system for classifying posts using categories, tags, and custom taxonomies.

### Architecture

- **Taxonomy** - A classification type (e.g., "category", "tag")
- **Term** - An individual item within a taxonomy (e.g., "Technology" category, "go" tag)
- **TermRelation** - Links posts to terms with optional ordering

### Key Features

- **Hierarchical terms** - Terms can have parent/child relationships for nested categories
- **Cached counts** - Each term stores a post count to avoid expensive queries
- **Slugs** - URL-friendly identifiers with uniqueness per taxonomy
- **Ordered relations** - Sequence field for manual term ordering on posts

### Usage Examples

```go
// Create a taxonomy
cat := blogstore.NewTaxonomy()
cat.SetName("Categories").SetSlug("category")
store.TaxonomyCreate(ctx, cat)

// Create hierarchical terms
tech := blogstore.NewTerm()
tech.SetTaxonomyID(cat.GetID()).SetName("Technology").SetSlug("technology")
store.TermCreate(ctx, tech)

prog := blogstore.NewTerm()
prog.SetTaxonomyID(cat.GetID()).
    SetParentID(tech.GetID()).
    SetName("Programming").
    SetSlug("programming")
store.TermCreate(ctx, prog)

// Assign terms to a post
store.PostSetTerms(ctx, postID, "category", []string{prog.GetID()})

// Get terms for a post
terms, _ := store.TermListByPostID(ctx, postID, "category")
```

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.
