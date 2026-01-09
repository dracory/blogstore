package mcp

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/blogstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestPostContentTypeMethods(t *testing.T) {
	// Test the new SetContentType and ContentType methods
	post := blogstore.NewPost()

	// Test default content type
	assert.Equal(t, "", post.ContentType(), "Default content type should be empty")
	assert.False(t, post.IsContentMarkdown(), "Should not be markdown by default")
	assert.False(t, post.IsContentHtml(), "Should not be HTML by default")
	assert.False(t, post.IsContentPlainText(), "Should not be plain text by default")

	// Test SetContentType with markdown
	post.SetContentType(blogstore.POST_CONTENT_TYPE_MARKDOWN)
	assert.Equal(t, blogstore.POST_CONTENT_TYPE_MARKDOWN, post.ContentType())
	assert.True(t, post.IsContentMarkdown(), "Should be markdown")
	assert.False(t, post.IsContentHtml(), "Should not be HTML")
	assert.False(t, post.IsContentPlainText(), "Should not be plain text")

	// Test SetContentType with HTML
	post.SetContentType(blogstore.POST_CONTENT_TYPE_HTML)
	assert.Equal(t, blogstore.POST_CONTENT_TYPE_HTML, post.ContentType())
	assert.False(t, post.IsContentMarkdown(), "Should not be markdown")
	assert.True(t, post.IsContentHtml(), "Should be HTML")
	assert.False(t, post.IsContentPlainText(), "Should not be plain text")

	// Test SetContentType with plain text
	post.SetContentType(blogstore.POST_CONTENT_TYPE_PLAIN_TEXT)
	assert.Equal(t, blogstore.POST_CONTENT_TYPE_PLAIN_TEXT, post.ContentType())
	assert.False(t, post.IsContentMarkdown(), "Should not be markdown")
	assert.False(t, post.IsContentHtml(), "Should not be HTML")
	assert.True(t, post.IsContentPlainText(), "Should be plain text")

	// Test that it's stored in metas
	assert.Equal(t, blogstore.POST_CONTENT_TYPE_PLAIN_TEXT, post.Meta("content_type"))
}

func TestContentTypeToEditor(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    string
	}{
		{
			name:        "markdown to Markdown",
			contentType: blogstore.POST_CONTENT_TYPE_MARKDOWN,
			expected:    blogstore.POST_EDITOR_MARKDOWN,
		},
		{
			name:        "html to HtmlArea",
			contentType: blogstore.POST_CONTENT_TYPE_HTML,
			expected:    blogstore.POST_EDITOR_HTMLAREA,
		},
		{
			name:        "plain_text to TextArea",
			contentType: blogstore.POST_CONTENT_TYPE_PLAIN_TEXT,
			expected:    blogstore.POST_EDITOR_TEXTAREA,
		},
		{
			name:        "unknown to TextArea",
			contentType: "unknown",
			expected:    blogstore.POST_EDITOR_TEXTAREA,
		},
		{
			name:        "empty to TextArea",
			contentType: "",
			expected:    blogstore.POST_EDITOR_TEXTAREA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contentTypeToEditor(tt.contentType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPostCreateWithContentType(t *testing.T) {
	// Setup
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store, err := blogstore.NewStore(blogstore.NewStoreOptions{
		PostTableName:      "test_posts",
		DB:                 db,
		DbDriverName:       "sqlite",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	mcp := NewMCP(store)

	// Test cases
	tests := []struct {
		name        string
		contentType string
		content     string
		expected    string
	}{
		{
			name:        "markdown content",
			contentType: "markdown",
			content:     "# Header\n**Bold text**",
			expected:    "markdown",
		},
		{
			name:        "html content",
			contentType: "html",
			content:     "<h1>Header</h1><strong>Bold text</strong>",
			expected:    "html",
		},
		{
			name:        "plain text content",
			contentType: "plain_text",
			content:     "Just plain text",
			expected:    "plain_text",
		},
		{
			name:        "default content_type",
			contentType: "",
			content:     "Some content",
			expected:    "plain_text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			request := map[string]any{
				"jsonrpc": "2.0",
				"id":      "1",
				"method":  "tools/call",
				"params": map[string]any{
					"name": "post_upsert",
					"arguments": map[string]any{
						"title":        "Test Post",
						"content":      tt.content,
						"content_type": tt.contentType,
						"status":       "draft",
					},
				},
			}

			reqBody, _ := json.Marshal(request)
			req := httptest.NewRequest("POST", "/mcp", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			// Execute
			mcp.Handler(w, req)

			// Check response
			assert.Equal(t, http.StatusOK, w.Code)

			var response struct {
				JSONRPC string `json:"jsonrpc"`
				ID      string `json:"id"`
				Result  struct {
					Content []struct {
						Type string `json:"type"`
						Text string `json:"text"`
					} `json:"content"`
				} `json:"result"`
			}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Parse the created post ID
			var createResult map[string]any
			err = json.Unmarshal([]byte(response.Result.Content[0].Text), &createResult)
			require.NoError(t, err)

			postID, ok := createResult["id"].(string)
			require.True(t, ok, "Post ID should be a string")

			// Retrieve the post to verify content_type was stored
			post, err := store.PostFindByID(context.Background(), postID)
			require.NoError(t, err)
			require.NotNil(t, post)

			// Check content_type using the new method
			storedContentType := post.ContentType()
			assert.Equal(t, tt.expected, storedContentType)

			// Check editor was set correctly
			expectedEditor := contentTypeToEditor(tt.expected)
			assert.Equal(t, expectedEditor, post.Editor())
		})
	}
}

func TestPostUpdateWithContentType(t *testing.T) {
	// Setup
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store, err := blogstore.NewStore(blogstore.NewStoreOptions{
		PostTableName:      "test_posts",
		DB:                 db,
		DbDriverName:       "sqlite",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	mcp := NewMCP(store)

	// Create initial post
	post := blogstore.NewPost()
	post.SetTitle("Initial Post")
	post.SetContent("Initial content")
	post.SetContentType(blogstore.POST_CONTENT_TYPE_PLAIN_TEXT)
	post.SetEditor(blogstore.POST_EDITOR_TEXTAREA)

	err = store.PostCreate(context.Background(), post)
	require.NoError(t, err)

	// Verify post was created and has an ID
	require.NotEmpty(t, post.ID(), "Post should have an ID after creation")

	// Test updating content_type
	tests := []struct {
		name        string
		contentType string
		content     string
		expected    string
	}{
		{
			name:        "update to markdown",
			contentType: "markdown",
			content:     "# Updated Header\n**Updated bold**",
			expected:    "markdown",
		},
		{
			name:        "update to html",
			contentType: "html",
			content:     "<h1>Updated Header</h1><strong>Updated bold</strong>",
			expected:    "html",
		},
		{
			name:        "update to plain text",
			contentType: "plain_text",
			content:     "Updated plain text",
			expected:    "plain_text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify post exists before update
			existingPost, err := store.PostFindByID(context.Background(), post.ID())
			require.NoError(t, err)
			require.NotNil(t, existingPost, "Post should exist before update")
			t.Logf("Updating post with ID: %s", post.ID())

			// Update request using post_upsert
			request := map[string]any{
				"jsonrpc": "2.0",
				"id":      "1",
				"method":  "tools/call",
				"params": map[string]any{
					"name": "post_upsert",
					"arguments": map[string]any{
						"id":           post.ID(),
						"content":      tt.content,
						"content_type": tt.contentType,
					},
				},
			}

			reqBody, _ := json.Marshal(request)
			req := httptest.NewRequest("POST", "/mcp", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			// Execute
			mcp.Handler(w, req)

			// Check response
			assert.Equal(t, http.StatusOK, w.Code)

			// Retrieve updated post
			updatedPost, err := store.PostFindByID(context.Background(), post.ID())
			require.NoError(t, err)
			require.NotNil(t, updatedPost)

			// Check content_type was updated using the new method
			storedContentType := updatedPost.ContentType()
			assert.Equal(t, tt.expected, storedContentType)

			// Check editor was updated correctly
			expectedEditor := contentTypeToEditor(tt.expected)
			assert.Equal(t, expectedEditor, updatedPost.Editor())

			// Check content was updated
			assert.Equal(t, tt.content, updatedPost.Content())
		})
	}
}

func TestPostCreateWithoutContentType(t *testing.T) {
	// Setup
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store, err := blogstore.NewStore(blogstore.NewStoreOptions{
		PostTableName:      "test_posts",
		DB:                 db,
		DbDriverName:       "sqlite",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	mcp := NewMCP(store)

	// Create request without content_type
	request := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_upsert",
			"arguments": map[string]any{
				"title":   "Test Post",
				"content": "Some content",
				"status":  "draft",
			},
		},
	}

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/mcp", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Execute
	mcp.Handler(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		JSONRPC string `json:"jsonrpc"`
		ID      string `json:"id"`
		Result  struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Parse the created post ID
	var createResult map[string]any
	err = json.Unmarshal([]byte(response.Result.Content[0].Text), &createResult)
	require.NoError(t, err)

	postID, ok := createResult["id"].(string)
	require.True(t, ok, "Post ID should be a string")

	// Retrieve the post to verify default content_type
	createdPost, err := store.PostFindByID(context.Background(), postID)
	require.NoError(t, err)
	require.NotNil(t, createdPost)

	// Check default content_type was set using the new method
	storedContentType := createdPost.ContentType()
	assert.Equal(t, blogstore.POST_CONTENT_TYPE_PLAIN_TEXT, storedContentType)

	// Check default editor was set
	assert.Equal(t, blogstore.POST_EDITOR_TEXTAREA, createdPost.Editor())
}

func TestBlogSchemaIncludesContentType(t *testing.T) {
	// Setup
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store, err := blogstore.NewStore(blogstore.NewStoreOptions{
		PostTableName:      "test_posts",
		DB:                 db,
		DbDriverName:       "sqlite",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	mcp := NewMCP(store)

	// Create request for blog schema
	request := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "blog_schema",
		},
	}

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/mcp", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Execute
	mcp.Handler(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		JSONRPC string `json:"jsonrpc"`
		ID      string `json:"id"`
		Result  struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Parse schema
	var schema map[string]any
	err = json.Unmarshal([]byte(response.Result.Content[0].Text), &schema)
	require.NoError(t, err)

	// Check content_type field is in schema
	entities, ok := schema["entities"].(map[string]any)
	require.True(t, ok)

	postEntity, ok := entities["post"].(map[string]any)
	require.True(t, ok)

	fields, ok := postEntity["fields"].([]any)
	require.True(t, ok)

	// Find content_type field
	var contentTypeField map[string]any
	found := false
	for _, field := range fields {
		fieldMap, ok := field.(map[string]any)
		require.True(t, ok)

		if name, ok := fieldMap["name"].(string); ok && name == "content_type" {
			contentTypeField = fieldMap
			found = true
			break
		}
	}

	require.True(t, found, "content_type field should be in schema")
	assert.Equal(t, "string", contentTypeField["type"])
	// Check enum values - handle type conversion from JSON unmarshaling
	enumInterface := contentTypeField["enum"]
	enumSlice := make([]string, 0)
	if enumList, ok := enumInterface.([]interface{}); ok {
		for _, item := range enumList {
			if str, ok := item.(string); ok {
				enumSlice = append(enumSlice, str)
			}
		}
	}
	assert.Equal(t, []string{"markdown", "html", "plain_text"}, enumSlice)
	assert.Contains(t, contentTypeField["description"], "Content format type")
}

func TestContentTypeValidation(t *testing.T) {
	// Setup
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store, err := blogstore.NewStore(blogstore.NewStoreOptions{
		PostTableName:      "test_posts",
		DB:                 db,
		DbDriverName:       "sqlite",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	mcp := NewMCP(store)

	// Test invalid content_type in schema validation
	request := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/list",
		"params":  map[string]any{},
	}

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/mcp", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Execute
	mcp.Handler(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		JSONRPC string `json:"jsonrpc"`
		ID      string `json:"id"`
		Result  struct {
			Tools []struct {
				Name        string         `json:"name"`
				InputSchema map[string]any `json:"inputSchema"`
			} `json:"tools"`
		} `json:"result"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Find post_upsert tool and check content_type validation
	var postUpsertTool struct {
		Name        string         `json:"name"`
		InputSchema map[string]any `json:"inputSchema"`
	}
	found := false
	for _, tool := range response.Result.Tools {
		if tool.Name == "post_upsert" {
			postUpsertTool = tool
			found = true
			break
		}
	}

	require.True(t, found, "post_upsert tool should be found")

	properties, ok := postUpsertTool.InputSchema["properties"].(map[string]any)
	require.True(t, ok)

	contentType, ok := properties["content_type"].(map[string]any)
	require.True(t, ok)

	// Check enum values - handle type conversion from JSON unmarshaling
	enumInterface := contentType["enum"]
	enumSlice := make([]string, 0)
	if enumList, ok := enumInterface.([]interface{}); ok {
		for _, item := range enumList {
			if str, ok := item.(string); ok {
				enumSlice = append(enumSlice, str)
			}
		}
	}
	require.True(t, len(enumSlice) > 0, "enum should have values")
	assert.Equal(t, []string{"markdown", "html", "plain_text"}, enumSlice)

	// Check default value
	defaultValue, ok := contentType["default"].(string)
	require.True(t, ok)
	assert.Equal(t, "plain_text", defaultValue)
}
