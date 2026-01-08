package mcp_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
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
	if !strings.Contains(respStr, "post_create") {
		t.Fatalf("Expected tools list to contain post_create: %s", respStr)
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

	// Create
	createReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_create",
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

	// Update (flat args)
	updateReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "3",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_update",
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

func Test_MCP_PostUpdate_WithNumericID_UsesNumber(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	ctx := context.Background()

	post := blogstore.NewPost().SetTitle("Numeric ID Test")
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate() error: %v", err)
	}

	// Send id as json.Number (simulates LLM numeric conversions without float64)
	updateReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "post_update",
			"arguments": map[string]any{
				"id":    json.Number(post.ID()),
				"title": "Updated",
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
	if !strings.Contains(updateText, "Updated") {
		t.Fatalf("Expected update response to contain updated title. Got: %s", updateText)
	}
}
