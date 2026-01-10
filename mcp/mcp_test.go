package mcp_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/blogstore"
	"github.com/dracory/blogstore/mcp"
	_ "modernc.org/sqlite"
)

func rpcResultText(t *testing.T, respBytes []byte) string {
	t.Helper()

	var rpcResp map[string]any
	if err := json.Unmarshal(respBytes, &rpcResp); err != nil {
		t.Fatalf("Failed to unmarshal json-rpc response: %v. Body=%s", err, string(respBytes))
	}

	result, ok := rpcResp["result"].(map[string]any)
	if !ok {
		t.Fatalf("Expected response to have result: %s", string(respBytes))
	}

	content, ok := result["content"].([]any)
	if !ok || len(content) == 0 {
		t.Fatalf("Expected response result.content: %s", string(respBytes))
	}

	item0, ok := content[0].(map[string]any)
	if !ok {
		t.Fatalf("Expected response result.content[0] object: %s", string(respBytes))
	}

	text, ok := item0["text"].(string)
	if !ok {
		t.Fatalf("Expected response result.content[0].text: %s", string(respBytes))
	}

	return text
}

func initDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := ":memory:?parseTime=true"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	return db
}

func initMCPServerWithStore(t *testing.T) (*httptest.Server, blogstore.StoreInterface, func()) {
	t.Helper()

	db := initDB(t)

	store, err := blogstore.NewStore(blogstore.NewStoreOptions{
		PostTableName:      "blog_posts",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	h := mcp.NewMCP(store)
	server := httptest.NewServer(http.HandlerFunc(h.Handler))
	return server, store, server.Close
}

func Test_MCP_Initialize(t *testing.T) {
	server, _, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	reqPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2025-06-18",
			"clientInfo": map[string]any{
				"name":    "test",
				"version": "0.0.0",
			},
		},
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	respStr := string(bodyBytes)
	if !strings.Contains(respStr, "protocolVersion") {
		t.Fatalf("Unexpected response: %s", respStr)
	}
	if !strings.Contains(respStr, "blogstore") {
		t.Fatalf("Expected serverInfo name blogstore: %s", respStr)
	}
}

func Test_MCP_ToolsList_StandardMethod(t *testing.T) {
	server, _, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	reqPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/list",
		"params":  map[string]any{},
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	respStr := string(bodyBytes)
	if !strings.Contains(respStr, "post_upsert") {
		t.Fatalf("Expected tools list to contain post_upsert: %s", respStr)
	}
	if !strings.Contains(respStr, "post_list") {
		t.Fatalf("Expected tools list to contain post_list: %s", respStr)
	}
	if !strings.Contains(respStr, "inputSchema") {
		t.Fatalf("Expected tools list to contain inputSchema: %s", respStr)
	}
}

func Test_MCP_ToolsCall_StandardMethod_PostCRUD(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	ctx := context.Background()

	// Create using upsert (no ID provided)
	createReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_upsert",
			"arguments": map[string]any{
				"title":   "Hello",
				"content": "World",
				"status":  blogstore.POST_STATUS_DRAFT,
			},
		},
	}
	createBody, _ := json.Marshal(createReq)
	createResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatalf("Failed to send create request: %v", err)
	}
	createRespBytes, _ := io.ReadAll(createResp.Body)
	createResp.Body.Close()

	createText := rpcResultText(t, createRespBytes)
	if !strings.Contains(createText, "Hello") {
		t.Fatalf("Expected create response to contain title. Got: %s", createText)
	}

	// Ensure we can list and obtain the id by inspecting the store directly (authoritative)
	posts, err := store.PostList(ctx, blogstore.PostQueryOptions{Limit: 10})
	if err != nil {
		t.Fatalf("PostList() error: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("Expected 1 post, got %d", len(posts))
	}
	postID := posts[0].ID()
	if postID == "" {
		t.Fatalf("Expected created post to have non-empty ID")
	}

	// Get
	getReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "2",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_get",
			"arguments": map[string]any{
				"id": postID,
			},
		},
	}
	getBody, _ := json.Marshal(getReq)
	getResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(getBody))
	if err != nil {
		t.Fatalf("Failed to send get request: %v", err)
	}
	getRespBytes, _ := io.ReadAll(getResp.Body)
	getResp.Body.Close()
	getText := rpcResultText(t, getRespBytes)
	if !strings.Contains(getText, postID) {
		t.Fatalf("Expected get response to contain id. Got: %s", getText)
	}

	// Update using upsert (ID provided)
	updateReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "3",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_upsert",
			"arguments": map[string]any{
				"id":    postID,
				"title": "Updated Title",
			},
		},
	}
	updateBody, _ := json.Marshal(updateReq)
	updateResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(updateBody))
	if err != nil {
		t.Fatalf("Failed to send update request: %v", err)
	}
	updateRespBytes, _ := io.ReadAll(updateResp.Body)
	updateResp.Body.Close()
	updateText := rpcResultText(t, updateRespBytes)
	if !strings.Contains(updateText, "Updated Title") {
		t.Fatalf("Expected update response to contain updated title. Got: %s", updateText)
	}

	// Delete
	deleteReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "4",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_delete",
			"arguments": map[string]any{
				"id": postID,
			},
		},
	}
	deleteBody, _ := json.Marshal(deleteReq)
	deleteResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(deleteBody))
	if err != nil {
		t.Fatalf("Failed to send delete request: %v", err)
	}
	deleteRespBytes, _ := io.ReadAll(deleteResp.Body)
	deleteResp.Body.Close()
	deleteText := rpcResultText(t, deleteRespBytes)
	if !strings.Contains(deleteText, "deleted") {
		t.Fatalf("Expected delete response to contain deleted flag. Got: %s", deleteText)
	}

	// Verify gone
	found, err := store.PostFindByID(ctx, postID)
	if err != nil {
		t.Fatalf("PostFindByID() error after delete: %v", err)
	}
	if found != nil {
		t.Fatalf("Expected post to be deleted")
	}
}

func Test_MCP_PostUpsert_CreateAndUpdate(t *testing.T) {
	db := initDB(t)
	defer db.Close()

	store, err := blogstore.NewStore(blogstore.NewStoreOptions{
		DB:                  db,
		PostTableName:       "posts",
		AutomigrateEnabled:  true,
		VersioningEnabled:   true,
		VersioningTableName: "versioning_table",
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	mcpServer := mcp.NewMCP(store)
	server := httptest.NewServer(http.HandlerFunc(mcpServer.Handler))
	defer server.Close()

	ctx := context.Background()

	// Test 1: Create new post with upsert (no ID provided)
	createReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_upsert",
			"arguments": map[string]any{
				"title":        "New Upsert Post",
				"content":      "Test content",
				"content_type": "markdown",
				"status":       "draft",
			},
		},
	}

	createBody, _ := json.Marshal(createReq)
	createResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatalf("Failed to send create upsert request: %v", err)
	}
	createRespBytes, _ := io.ReadAll(createResp.Body)
	createResp.Body.Close()

	createText := rpcResultText(t, createRespBytes)
	var createResult map[string]any
	if err := json.Unmarshal([]byte(createText), &createResult); err != nil {
		t.Fatalf("Failed to parse create result: %v", err)
	}

	postID, ok := createResult["id"].(string)
	if !ok || postID == "" {
		t.Fatalf("Expected create response to contain post ID. Got: %s", createText)
	}

	// Verify post was created
	post, err := store.PostFindByID(ctx, postID)
	if err != nil {
		t.Fatalf("Failed to find created post: %v", err)
	}
	if post == nil {
		t.Fatalf("Created post not found")
	}

	// Check that no versions exist yet (only creation, no updates)
	versions, err := store.VersioningList(ctx, blogstore.NewVersioningQuery().
		SetEntityType(blogstore.VERSIONING_TYPE_POST).
		SetEntityID(postID))
	if err != nil {
		t.Fatalf("Failed to list versions: %v", err)
	}
	initialVersionCount := len(versions)
	if initialVersionCount != 1 {
		t.Fatalf("Expected 1 version after create, got %d", initialVersionCount)
	}
	t.Logf("Initial version count: %d", initialVersionCount)

	// Test 2: Update existing post with upsert (ID provided)
	updateReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "2",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_upsert",
			"arguments": map[string]any{
				"id":           postID,
				"title":        "Updated Upsert Post",
				"content":      "Updated content",
				"content_type": "markdown",
				"status":       "published",
				"featured":     "yes",
			},
		},
	}

	updateBody, _ := json.Marshal(updateReq)
	updateResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(updateBody))
	if err != nil {
		t.Fatalf("Failed to send update upsert request: %v", err)
	}
	updateRespBytes, _ := io.ReadAll(updateResp.Body)
	updateResp.Body.Close()

	updateText := rpcResultText(t, updateRespBytes)
	var updateResult map[string]any
	if err := json.Unmarshal([]byte(updateText), &updateResult); err != nil {
		t.Fatalf("Failed to parse update result: %v", err)
	}

	if updateResult["id"].(string) != postID {
		t.Fatalf("Expected same post ID after update. Got: %v", updateResult["id"])
	}

	if updateResult["title"].(string) != "Updated Upsert Post" {
		t.Fatalf("Expected updated title. Got: %v", updateResult["title"])
	}

	// Verify post was updated
	updatedPost, err := store.PostFindByID(ctx, postID)
	if err != nil {
		t.Fatalf("Failed to find updated post: %v", err)
	}
	if updatedPost.Title() != "Updated Upsert Post" {
		t.Fatalf("Expected post title to be updated. Got: %s", updatedPost.Title())
	}
	if updatedPost.Status() != "published" {
		t.Fatalf("Expected post status to be published. Got: %s", updatedPost.Status())
	}
	if updatedPost.Featured() != "yes" {
		t.Fatalf("Expected post to be featured. Got: %s", updatedPost.Featured())
	}

	// Check that a new version was created
	versionsAfterUpdate, err := store.VersioningList(ctx, blogstore.NewVersioningQuery().
		SetEntityType(blogstore.VERSIONING_TYPE_POST).
		SetEntityID(postID).
		SetOrderBy("created_at").
		SetSortOrder("DESC"))
	if err != nil {
		t.Fatalf("Failed to list versions after update: %v", err)
	}

	expectedVersionCount := initialVersionCount + 1
	if len(versionsAfterUpdate) != expectedVersionCount {
		t.Fatalf("Expected %d versions after update, got %d", expectedVersionCount, len(versionsAfterUpdate))
	}

	// Verify the latest version contains the post-update state
	if len(versionsAfterUpdate) > 0 {
		foundUpdated := false
		for _, v := range versionsAfterUpdate {
			versionContent := v.Content()
			var versionedPostData map[string]interface{}
			if err := json.Unmarshal([]byte(versionContent), &versionedPostData); err != nil {
				t.Fatalf("Failed to parse version content: %v", err)
			}
			if versionedPostData["title"] == "Updated Upsert Post" {
				foundUpdated = true
				break
			}
		}

		if !foundUpdated {
			t.Fatalf("Expected versions to contain updated title 'Updated Upsert Post'")
		}
	}
}

func Test_MCP_PostUpsert_VersioningIntegration(t *testing.T) {
	db := initDB(t)
	defer db.Close()

	store, err := blogstore.NewStore(blogstore.NewStoreOptions{
		DB:                  db,
		PostTableName:       "posts",
		AutomigrateEnabled:  true,
		VersioningEnabled:   true,
		VersioningTableName: "versioning_table",
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	mcpServer := mcp.NewMCP(store)
	server := httptest.NewServer(http.HandlerFunc(mcpServer.Handler))
	defer server.Close()

	ctx := context.Background()

	// Create initial post
	createReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_upsert",
			"arguments": map[string]any{
				"title":   "Version Test Post",
				"content": "Initial content",
				"status":  "draft",
			},
		},
	}

	createBody, _ := json.Marshal(createReq)
	createResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatalf("Failed to send create request: %v", err)
	}
	createRespBytes, _ := io.ReadAll(createResp.Body)
	createResp.Body.Close()

	createText := rpcResultText(t, createRespBytes)
	var createResult map[string]any
	if err := json.Unmarshal([]byte(createText), &createResult); err != nil {
		t.Fatalf("Failed to parse create result: %v", err)
	}

	postID := createResult["id"].(string)

	// Perform multiple updates to create multiple versions
	updates := []struct {
		title   string
		content string
		status  string
	}{
		{"First Update", "First updated content", "draft"},
		{"Second Update", "Second updated content", "published"},
		{"Third Update", "Third updated content", "published"},
	}

	for i, update := range updates {
		updateReq := map[string]any{
			"jsonrpc": "2.0",
			"id":      fmt.Sprintf("%d", i+2),
			"method":  "tools/call",
			"params": map[string]any{
				"name": "post_upsert",
				"arguments": map[string]any{
					"id":      postID,
					"title":   update.title,
					"content": update.content,
					"status":  update.status,
				},
			},
		}

		updateBody, _ := json.Marshal(updateReq)
		updateResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(updateBody))
		if err != nil {
			t.Fatalf("Failed to send update request %d: %v", i+1, err)
		}
		updateRespBytes, _ := io.ReadAll(updateResp.Body)
		updateResp.Body.Close()

		// Verify update succeeded
		updateText := rpcResultText(t, updateRespBytes)
		var updateResult map[string]any
		if err := json.Unmarshal([]byte(updateText), &updateResult); err != nil {
			t.Fatalf("Failed to parse update result %d: %v", i+1, err)
		}

		if updateResult["title"].(string) != update.title {
			t.Fatalf("Update %d: Expected title '%s', got '%s'", i+1, update.title, updateResult["title"])
		}
	}

	// Verify we have the correct number of versions
	versions, err := store.VersioningList(ctx, blogstore.NewVersioningQuery().
		SetEntityType(blogstore.VERSIONING_TYPE_POST).
		SetEntityID(postID).
		SetOrderBy("created_at").
		SetSortOrder("DESC"))
	if err != nil {
		t.Fatalf("Failed to list versions: %v", err)
	}

	expectedVersionCount := 1 + len(updates) // One on create + one per update
	if len(versions) != expectedVersionCount {
		t.Fatalf("Expected %d versions, got %d", expectedVersionCount, len(versions))
	}

	// Verify version history contains the expected titles. Do not assume ordering
	// because multiple versions may share the same created_at value.
	expectedTitles := map[string]bool{
		"Version Test Post": false,
		"First Update":      false,
		"Second Update":     false,
		"Third Update":      false,
	}

	for i, version := range versions {
		versionContent := version.Content()
		var versionedPostData map[string]interface{}
		if err := json.Unmarshal([]byte(versionContent), &versionedPostData); err != nil {
			t.Fatalf("Failed to parse version %d content: %v", i, err)
		}
		if title, ok := versionedPostData["title"].(string); ok {
			if _, exists := expectedTitles[title]; exists {
				expectedTitles[title] = true
			}
		}
	}

	for title, found := range expectedTitles {
		if !found {
			t.Fatalf("Expected versions to contain title '%s'", title)
		}
	}

	t.Logf("Successfully verified %d versions for post %s", len(versions), postID)
}

func Test_MCP_PostVersions(t *testing.T) {
	db := initDB(t)
	defer db.Close()

	store, err := blogstore.NewStore(blogstore.NewStoreOptions{
		DB:                  db,
		PostTableName:       "posts",
		AutomigrateEnabled:  true,
		VersioningEnabled:   true,
		VersioningTableName: "versioning_table",
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	mcpServer := mcp.NewMCP(store)
	server := httptest.NewServer(http.HandlerFunc(mcpServer.Handler))
	defer server.Close()

	// Create initial post
	createReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_upsert",
			"arguments": map[string]any{
				"title":   "Versions Test Post",
				"content": "Initial content",
				"status":  "draft",
			},
		},
	}

	createBody, _ := json.Marshal(createReq)
	createResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatalf("Failed to send create request: %v", err)
	}
	createRespBytes, _ := io.ReadAll(createResp.Body)
	createResp.Body.Close()

	createText := rpcResultText(t, createRespBytes)
	var createResult map[string]any
	if err := json.Unmarshal([]byte(createText), &createResult); err != nil {
		t.Fatalf("Failed to parse create result: %v", err)
	}

	postID := createResult["id"].(string)

	// Perform a few updates to create versions
	for i := 1; i <= 3; i++ {
		updateReq := map[string]any{
			"jsonrpc": "2.0",
			"id":      fmt.Sprintf("%d", i+1),
			"method":  "tools/call",
			"params": map[string]any{
				"name": "post_upsert",
				"arguments": map[string]any{
					"id":      postID,
					"title":   fmt.Sprintf("Update %d", i),
					"content": fmt.Sprintf("Content %d", i),
				},
			},
		}

		updateBody, _ := json.Marshal(updateReq)
		updateResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(updateBody))
		if err != nil {
			t.Fatalf("Failed to send update request %d: %v", i, err)
		}
		updateResp.Body.Close()
	}

	// Test post_versions tool
	versionsReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "5",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_versions",
			"arguments": map[string]any{
				"id": postID,
			},
		},
	}

	versionsBody, _ := json.Marshal(versionsReq)
	versionsResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(versionsBody))
	if err != nil {
		t.Fatalf("Failed to send versions request: %v", err)
	}
	versionsRespBytes, _ := io.ReadAll(versionsResp.Body)
	versionsResp.Body.Close()

	versionsText := rpcResultText(t, versionsRespBytes)
	var versionsResult map[string]any
	if err := json.Unmarshal([]byte(versionsText), &versionsResult); err != nil {
		t.Fatalf("Failed to parse versions result: %v", err)
	}

	// Verify versions response
	if versionsResult["total"].(float64) != 4 {
		t.Fatalf("Expected 4 versions, got: %v", versionsResult["total"])
	}

	versions, ok := versionsResult["versions"].([]any)
	if !ok {
		t.Fatalf("Expected versions array, got: %T", versionsResult["versions"])
	}

	if len(versions) != 4 {
		t.Fatalf("Expected 4 version items, got: %d", len(versions))
	}

	// Verify version structure
	for i, version := range versions {
		versionMap, ok := version.(map[string]any)
		if !ok {
			t.Fatalf("Version %d should be a map, got: %T", i, version)
		}

		// Check required fields
		requiredFields := []string{"id", "entity_id", "entity_type", "content", "created_at"}
		for _, field := range requiredFields {
			if _, exists := versionMap[field]; !exists {
				t.Fatalf("Version %d missing required field: %s", i, field)
			}
		}

		// Verify entity type and ID
		if versionMap["entity_type"] != "post" {
			t.Fatalf("Version %d: Expected entity_type 'post', got: %v", i, versionMap["entity_type"])
		}

		if versionMap["entity_id"] != postID {
			t.Fatalf("Version %d: Expected entity_id '%s', got: %v", i, postID, versionMap["entity_id"])
		}

		// Verify content is valid JSON (post data)
		content, ok := versionMap["content"].(string)
		if !ok {
			t.Fatalf("Version %d: Content should be a string, got: %T", i, versionMap["content"])
		}

		var postData map[string]interface{}
		if err := json.Unmarshal([]byte(content), &postData); err != nil {
			t.Fatalf("Version %d: Content should be valid JSON, got error: %v", i, err)
		}

		if postData["title"] == nil {
			t.Fatalf("Version %d: Content should contain title field", i)
		}

		t.Logf("Version %d: Validated with title '%v'", i, postData["title"])
	}

	t.Logf("Successfully validated post_versions tool for post %s", postID)
}
