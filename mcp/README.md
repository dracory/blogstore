# Blog Store MCP Handler

This package provides an MCP (Model Context Protocol) handler for Blog Store, enabling LLM clients to manage blog posts via JSON-RPC tools.

## Features

- Post management (create, list, read, update, delete)
- JSON-RPC 2.0 compatible
- Attachable to any existing HTTP server

## Getting Started

### Basic Usage

```go
package main

import (
	"log"
	"net/http"

	"github.com/dracory/blogstore"
	"github.com/dracory/blogstore/mcp"
)

func main() {
	// Initialize your blog store
	store, err := blogstore.NewStore(blogstore.NewStoreOptions{
		DB:                 db, // your *sql.DB
		PostTableName:      "blog_posts",
		AutomigrateEnabled: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	mcpHandler := mcp.NewMCP(store)

	http.HandleFunc("/mcp/blog", mcpHandler.Handler)
	log.Println("Starting server on :8080")
	_ = http.ListenAndServe(":8080", nil)
}
```

## MCP Protocol

This handler supports both MCP-standard JSON-RPC methods and legacy aliases:

- MCP-standard:
  - `initialize`
  - `notifications/initialized`
  - `tools/list`
  - `tools/call`
- Legacy aliases:
  - `list_tools` (alias of `tools/list`)
  - `call_tool` (alias of `tools/call`)

### List tools

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "tools/list",
  "params": {}
}
```

### Call a tool

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "tools/call",
  "params": {
    "name": "post_list",
    "arguments": {
      "limit": 10,
      "offset": 0
    }
  }
}
```

Tool results are returned as:

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{\"items\":[]}" 
      }
    ]
  }
}
```

## Supported Tools

- `blog_schema` - Get detailed schema information and field constraints
- `post_list` - List blog posts with filtering options
- `post_create` - Create a new blog post
- `post_get` - Get a blog post by ID
- `post_update` - Update an existing blog post
- `post_delete` - Delete a blog post

## Important Field Constraints

### Featured Field
The `featured` field requires string values, not boolean:
- ✅ **Correct**: `"featured": "yes"` or `"featured": "no"`
- ❌ **Incorrect**: `"featured": true` or `"featured": false`

### Status Field
The `status` field supports specific string values:
- `"draft"` - Default status, not publicly visible
- `"published"` - Publicly visible posts
- `"unpublished"` - Hidden posts
- `"trash"` - Deleted posts

### Content Format
The `content` field supports multiple formats:
- **Markdown** (recommended for blog posts)
- **HTML**
- **Plain text**

### Content Type Field
The `content_type` field specifies how content should be rendered:
- `"markdown"` - For Markdown content (enables proper rendering with headers, emphasis, etc.)
- `"plain_text"` - Default for plain text content
- `"html"` - For HTML content

**Important**: Set `content_type="markdown"` when creating markdown content to ensure proper rendering!

## Schema Discovery

Use the `blog_schema` tool to get detailed information about entities and their constraints:

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "tools/call",
  "params": {
    "name": "blog_schema",
    "arguments": {}
  }
}
```

This returns comprehensive schema information including:
- Field definitions with types and descriptions
- Allowed values for enum fields
- Usage notes and best practices
- Tool argument specifications

## ID typing: always use strings

Some clients/LLMs may convert large integer-looking strings into JSON numbers (sometimes scientific notation), which is lossy.

- Always send identifiers as strings:

```json
{ "id": "20260108160058473" }
```

The handler advertises identifier fields as `type: "string"` via `tools/list` `inputSchema` to help clients send the correct types.
