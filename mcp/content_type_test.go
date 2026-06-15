package mcp

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/dracory/blogstore"
	_ "modernc.org/sqlite"
)

func TestPostContentTypeMethods(t *testing.T) {
	// Test the new SetContentType and ContentType methods
	post := blogstore.NewPost()

	// Test default content type
	if post.GetContentType() != "" {
		t.Errorf("Default content type should be empty, got %q", post.GetContentType())
	}
	if post.IsContentMarkdown() {
		t.Errorf("Should not be markdown by default")
	}
	if post.IsContentHtml() {
		t.Errorf("Should not be HTML by default")
	}
	if post.IsContentPlainText() {
		t.Errorf("Should not be plain text by default")
	}

	// Test SetContentType with markdown
	post.SetContentType(blogstore.POST_CONTENT_TYPE_MARKDOWN)
	if post.GetContentType() != blogstore.POST_CONTENT_TYPE_MARKDOWN {
		t.Errorf("Expected content type to be %s, got %s", blogstore.POST_CONTENT_TYPE_MARKDOWN, post.GetContentType())
	}
	if !post.IsContentMarkdown() {
		t.Errorf("Should be markdown")
	}
	if post.IsContentHtml() {
		t.Errorf("Should not be HTML")
	}
	if post.IsContentPlainText() {
		t.Errorf("Should not be plain text")
	}

	// Test SetContentType with HTML
	post.SetContentType(blogstore.POST_CONTENT_TYPE_HTML)
	if post.GetContentType() != blogstore.POST_CONTENT_TYPE_HTML {
		t.Errorf("Expected content type to be %s, got %s", blogstore.POST_CONTENT_TYPE_HTML, post.GetContentType())
	}
	if post.IsContentMarkdown() {
		t.Errorf("Should not be markdown")
	}
	if !post.IsContentHtml() {
		t.Errorf("Should be HTML")
	}
	if post.IsContentPlainText() {
		t.Errorf("Should not be plain text")
	}

	// Test SetContentType with plain text
	post.SetContentType(blogstore.POST_CONTENT_TYPE_PLAIN_TEXT)
	if post.GetContentType() != blogstore.POST_CONTENT_TYPE_PLAIN_TEXT {
		t.Errorf("Expected content type to be %s, got %s", blogstore.POST_CONTENT_TYPE_PLAIN_TEXT, post.GetContentType())
	}
	if post.IsContentMarkdown() {
		t.Errorf("Should not be markdown")
	}
	if post.IsContentHtml() {
		t.Errorf("Should not be HTML")
	}
	if !post.IsContentPlainText() {
		t.Errorf("Should be plain text")
	}

	// Test that it's stored in metas
	if post.GetMeta("content_type") != blogstore.POST_CONTENT_TYPE_PLAIN_TEXT {
		t.Errorf("Expected content_type meta to be %s, got %s", blogstore.POST_CONTENT_TYPE_PLAIN_TEXT, post.GetMeta("content_type"))
	}
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
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
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
			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}

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
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// Parse the created post ID
			var createResult map[string]any
			err = json.Unmarshal([]byte(response.Result.Content[0].Text), &createResult)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			postID, ok := createResult["id"].(string)
			if !ok {
				t.Fatalf("Post ID should be a string")
			}

			// Retrieve the post to verify content_type was stored
			post, err := store.PostFindByID(context.Background(), postID)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if post == nil {
				t.Fatalf("Expected non-nil post")
			}

			// Check content_type using the new method
			storedContentType := post.GetContentType()
			if storedContentType != tt.expected {
				t.Errorf("Expected content type %s, got %s", tt.expected, storedContentType)
			}

			// Check editor was set correctly
			expectedEditor := contentTypeToEditor(tt.expected)
			if post.GetEditor() != expectedEditor {
				t.Errorf("Expected editor %s, got %s", expectedEditor, post.GetEditor())
			}
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
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify post was created and has an ID
	if post.GetID() == "" {
		t.Fatalf("Post should have an ID after creation")
	}

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
			existingPost, err := store.PostFindByID(context.Background(), post.GetID())
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if existingPost == nil {
				t.Fatalf("Post should exist before update")
			}
			t.Logf("Updating post with ID: %s", post.GetID())

			// Update request using post_upsert
			request := map[string]any{
				"jsonrpc": "2.0",
				"id":      "1",
				"method":  "tools/call",
				"params": map[string]any{
					"name": "post_upsert",
					"arguments": map[string]any{
						"id":           post.GetID(),
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
			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}

			// Retrieve updated post
			updatedPost, err := store.PostFindByID(context.Background(), post.GetID())
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if updatedPost == nil {
				t.Fatalf("Expected non-nil post")
			}

			// Check content_type was updated using the new method
			storedContentType := updatedPost.GetContentType()
			if storedContentType != tt.expected {
				t.Errorf("Expected content type %s, got %s", tt.expected, storedContentType)
			}

			// Check editor was updated correctly
			expectedEditor := contentTypeToEditor(tt.expected)
			if updatedPost.GetEditor() != expectedEditor {
				t.Errorf("Expected editor %s, got %s", expectedEditor, updatedPost.GetEditor())
			}

			// Check content was updated
			if updatedPost.GetContent() != tt.content {
				t.Errorf("Expected content %q, got %q", tt.content, updatedPost.GetContent())
			}
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
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

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
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Parse the created post ID
	var createResult map[string]any
	err = json.Unmarshal([]byte(response.Result.Content[0].Text), &createResult)
	if err != nil {
		t.Errorf("Failed to unmarshal create result: %v", err)
	}

	postID, ok := createResult["id"].(string)
	if !ok {
		t.Errorf("Post ID should be a string")
	}

	// Retrieve the post to verify default content_type
	createdPost, err := store.PostFindByID(context.Background(), postID)
	if err != nil {
		t.Errorf("Failed to retrieve post: %v", err)
	}
	if createdPost == nil {
		t.Errorf("Expected non-nil post")
	}

	// Check default content_type was set using the new method
	storedContentType := createdPost.GetContentType()
	if storedContentType != blogstore.POST_CONTENT_TYPE_PLAIN_TEXT {
		t.Errorf("Expected content type %s, got %s", blogstore.POST_CONTENT_TYPE_PLAIN_TEXT, storedContentType)
	}

	// Check default editor was set
	if createdPost.GetEditor() != blogstore.POST_EDITOR_TEXTAREA {
		t.Errorf("Expected editor %s, got %s", blogstore.POST_EDITOR_TEXTAREA, createdPost.GetEditor())
	}
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
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

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
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Parse schema
	var schema map[string]any
	err = json.Unmarshal([]byte(response.Result.Content[0].Text), &schema)
	if err != nil {
		t.Errorf("Failed to unmarshal schema: %v", err)
	}

	// Check content_type field is in schema
	entities, ok := schema["entities"].(map[string]any)
	if !ok {
		t.Fatalf("Expected entities to be a map")
	}

	postEntity, ok := entities["post"].(map[string]any)
	if !ok {
		t.Fatalf("Expected post to be a map")
	}

	fields, ok := postEntity["fields"].([]any)
	if !ok {
		t.Fatalf("Expected fields to be a slice")
	}

	// Find content_type field
	var contentTypeField map[string]any
	found := false
	for _, field := range fields {
		fieldMap, ok := field.(map[string]any)
		if !ok {
			t.Fatalf("Expected field to be a map")
		}

		if name, ok := fieldMap["name"].(string); ok && name == "content_type" {
			contentTypeField = fieldMap
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("content_type field should be in schema")
	}
	if contentTypeField["type"] != "string" {
		t.Errorf("Expected type to be string, got %v", contentTypeField["type"])
	}
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
	if !reflect.DeepEqual(enumSlice, []string{"markdown", "html", "plain_text"}) {
		t.Errorf("Expected enum to be [markdown, html, plain_text], got %v", enumSlice)
	}
	desc, ok := contentTypeField["description"].(string)
	if !ok || !strings.Contains(desc, "Content format type") {
		t.Errorf("Expected description to contain 'Content format type', got %v", contentTypeField["description"])
	}
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
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

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
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

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

	if !found {
		t.Fatalf("post_upsert tool should be found")
	}

	properties, ok := postUpsertTool.InputSchema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("Expected properties to be a map")
	}

	contentType, ok := properties["content_type"].(map[string]any)
	if !ok {
		t.Fatalf("Expected content_type to be a map")
	}

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
	if len(enumSlice) == 0 {
		t.Fatalf("enum should have values")
	}
	if !reflect.DeepEqual(enumSlice, []string{"markdown", "html", "plain_text"}) {
		t.Errorf("Expected enum to be [markdown, html, plain_text], got %v", enumSlice)
	}

	// Check default value
	defaultValue, ok := contentType["default"].(string)
	if !ok {
		t.Fatalf("Expected default to be a string")
	}
	if defaultValue != "plain_text" {
		t.Errorf("Expected default to be 'plain_text', got %s", defaultValue)
	}
}
