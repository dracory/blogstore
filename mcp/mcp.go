package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dracory/blogstore"
)

type MCP struct {
	store blogstore.StoreInterface
}

func NewMCP(store blogstore.StoreInterface) *MCP {
	return &MCP{store: store}
}

// Handler is an HTTP handler intended to be mounted at a dedicated route.
//
// The protocol is JSON-RPC 2.0 compatible and currently supports:
// - MCP standard methods: initialize, notifications/initialized, tools/list, tools/call
// - legacy aliases: list_tools, call_tool
func (m *MCP) Handler(w http.ResponseWriter, r *http.Request) {
	if m == nil || m.store == nil {
		writeJSON(w, http.StatusInternalServerError, jsonRPCErrorResponse(nil, -32603, "store is not initialized"))
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, jsonRPCErrorResponse(nil, -32602, "failed to read request body"))
		return
	}
	defer r.Body.Close()

	var req jsonRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSON(w, http.StatusOK, jsonRPCErrorResponse(nil, -32700, "parse error"))
		return
	}

	if strings.TrimSpace(req.JSONRPC) == "" {
		req.JSONRPC = "2.0"
	}

	switch req.Method {
	case "initialize":
		m.handleInitialize(w, r.Context(), req.ID, req.Params)
		return
	case "notifications/initialized":
		m.handleInitialized(w, r.Context())
		return
	case "tools/list":
		m.handleToolsList(w, r.Context(), req.ID)
		return
	case "tools/call":
		m.handleToolsCall(w, r.Context(), req.ID, req.Params)
		return
	case "list_tools":
		m.handleToolsList(w, r.Context(), req.ID)
		return
	case "call_tool":
		m.handleToolsCall(w, r.Context(), req.ID, req.Params)
		return
	default:
		writeJSON(w, http.StatusOK, jsonRPCErrorResponse(req.ID, -32601, "method not found"))
		return
	}
}

func argString(args map[string]any, key string) string {
	v, ok := args[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case json.Number:
		return t.String()
	case float64:
		return fmt.Sprintf("%.0f", t)
	case int:
		return fmt.Sprintf("%d", t)
	case int64:
		return fmt.Sprintf("%d", t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func argInt(args map[string]any, key string) (int, bool) {
	v, ok := args[key]
	if !ok || v == nil {
		return 0, false
	}
	switch t := v.(type) {
	case json.Number:
		i64, err := t.Int64()
		if err != nil {
			return 0, false
		}
		return int(i64), true
	case float64:
		return int(t), true
	case int:
		return t, true
	case int64:
		return int(t), true
	default:
		return 0, false
	}
}

func argBool(args map[string]any, key string) (bool, bool) {
	v, ok := args[key]
	if !ok || v == nil {
		return false, false
	}
	switch t := v.(type) {
	case bool:
		return t, true
	case string:
		vv := strings.TrimSpace(strings.ToLower(t))
		if vv == "true" || vv == "1" || vv == "yes" {
			return true, true
		}
		if vv == "false" || vv == "0" || vv == "no" {
			return false, true
		}
		return false, false
	default:
		return false, false
	}
}

type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type jsonRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func jsonRPCErrorResponse(id any, code int, message string) jsonRPCResponse {
	return jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: jsonRPCError{
			Code:    code,
			Message: message,
		},
	}
}

func jsonRPCResultResponse(id any, result any) jsonRPCResponse {
	return jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func toolTextResult(text string) map[string]any {
	return map[string]any{
		"content": []map[string]any{
			{
				"type": "text",
				"text": text,
			},
		},
	}
}

func (m *MCP) handleInitialize(w http.ResponseWriter, ctx context.Context, id any, params json.RawMessage) {
	var p struct {
		ProtocolVersion string `json:"protocolVersion"`
		ClientInfo      any    `json:"clientInfo"`
		Capabilities    any    `json:"capabilities"`
	}
	_ = json.Unmarshal(params, &p)

	result := map[string]any{
		"protocolVersion": "2025-06-18",
		"serverInfo": map[string]any{
			"name":    "blogstore",
			"version": "0.1.0",
		},
		"capabilities": map[string]any{
			"tools": map[string]any{},
		},
		"echo": map[string]any{
			"clientProtocolVersion": p.ProtocolVersion,
			"clientInfo":            p.ClientInfo,
			"clientCapabilities":    p.Capabilities,
		},
	}

	writeJSON(w, http.StatusOK, jsonRPCResultResponse(id, result))
}

func (m *MCP) handleInitialized(w http.ResponseWriter, ctx context.Context) {
	w.WriteHeader(http.StatusOK)
}

func (m *MCP) handleToolsList(w http.ResponseWriter, ctx context.Context, id any) {
	tools := []map[string]any{
		{
			"name":        "post_list",
			"description": "List blog posts",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit":        map[string]any{"type": "integer"},
					"offset":       map[string]any{"type": "integer"},
					"id":           map[string]any{"type": "string"},
					"status":       map[string]any{"type": "string"},
					"search":       map[string]any{"type": "string"},
					"with_deleted": map[string]any{"type": "boolean"},
					"order_by":     map[string]any{"type": "string"},
					"sort_order":   map[string]any{"type": "string"},
				},
			},
		},
		{
			"name":        "post_create",
			"description": "Create a blog post",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"title"},
				"properties": map[string]any{
					"id":               map[string]any{"type": "string"},
					"title":            map[string]any{"type": "string"},
					"content":          map[string]any{"type": "string"},
					"summary":          map[string]any{"type": "string"},
					"status":           map[string]any{"type": "string"},
					"author_id":        map[string]any{"type": "string"},
					"canonical_url":    map[string]any{"type": "string"},
					"image_url":        map[string]any{"type": "string"},
					"featured":         map[string]any{"type": "string"},
					"published_at":     map[string]any{"type": "string"},
					"meta_description": map[string]any{"type": "string"},
					"meta_keywords":    map[string]any{"type": "string"},
					"meta_robots":      map[string]any{"type": "string"},
					"memo":             map[string]any{"type": "string"},
				},
			},
		},
		{
			"name":        "post_get",
			"description": "Get a blog post by ID",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"id"},
				"properties": map[string]any{
					"id": map[string]any{"type": "string"},
				},
			},
		},
		{
			"name":        "post_update",
			"description": "Update a blog post",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"id"},
				"properties": map[string]any{
					"id": map[string]any{"type": "string"},
					"updates": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"title":            map[string]any{"type": "string"},
							"content":          map[string]any{"type": "string"},
							"summary":          map[string]any{"type": "string"},
							"status":           map[string]any{"type": "string"},
							"author_id":        map[string]any{"type": "string"},
							"canonical_url":    map[string]any{"type": "string"},
							"image_url":        map[string]any{"type": "string"},
							"featured":         map[string]any{"type": "string"},
							"published_at":     map[string]any{"type": "string"},
							"meta_description": map[string]any{"type": "string"},
							"meta_keywords":    map[string]any{"type": "string"},
							"meta_robots":      map[string]any{"type": "string"},
							"memo":             map[string]any{"type": "string"},
						},
					},
				},
			},
		},
		{
			"name":        "post_delete",
			"description": "Delete a blog post",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"id"},
				"properties": map[string]any{
					"id": map[string]any{"type": "string"},
				},
			},
		},
	}

	result := map[string]any{"tools": tools}
	writeJSON(w, http.StatusOK, jsonRPCResultResponse(id, result))
}

func (m *MCP) handleToolsCall(w http.ResponseWriter, ctx context.Context, id any, params json.RawMessage) {
	var p struct {
		Name      string          `json:"name"`
		ToolName  string          `json:"tool_name"`
		Args      json.RawMessage `json:"arguments"`
		Arguments json.RawMessage `json:"params"`
	}
	_ = json.Unmarshal(params, &p)

	toolName := strings.TrimSpace(p.Name)
	if toolName == "" {
		toolName = strings.TrimSpace(p.ToolName)
	}

	argsRaw := p.Args
	if len(argsRaw) == 0 {
		argsRaw = p.Arguments
	}

	args := map[string]any{}
	if len(argsRaw) > 0 {
		dec := json.NewDecoder(strings.NewReader(string(argsRaw)))
		dec.UseNumber()
		if err := dec.Decode(&args); err != nil {
			writeJSON(w, http.StatusOK, jsonRPCErrorResponse(id, -32602, "invalid tool arguments"))
			return
		}
	}

	text, err := m.dispatchTool(ctx, toolName, args)
	if err != nil {
		writeJSON(w, http.StatusOK, jsonRPCErrorResponse(id, -32603, err.Error()))
		return
	}

	writeJSON(w, http.StatusOK, jsonRPCResultResponse(id, toolTextResult(text)))
}

func (m *MCP) dispatchTool(ctx context.Context, toolName string, args map[string]any) (string, error) {
	switch toolName {
	case "post_list":
		return m.toolPostList(ctx, args)
	case "post_create":
		return m.toolPostCreate(ctx, args)
	case "post_get":
		return m.toolPostGet(ctx, args)
	case "post_update":
		return m.toolPostUpdate(ctx, args)
	case "post_delete":
		return m.toolPostDelete(ctx, args)
	default:
		return "", errors.New("unknown tool")
	}
}

func postToMap(post *blogstore.Post) map[string]string {
	if post == nil {
		return map[string]string{}
	}
	return post.Data()
}

func (m *MCP) toolPostList(ctx context.Context, args map[string]any) (string, error) {
	opts := blogstore.PostQueryOptions{}

	opts.ID = argString(args, "id")
	opts.Status = argString(args, "status")
	opts.Search = argString(args, "search")
	opts.OrderBy = argString(args, "order_by")
	opts.SortOrder = argString(args, "sort_order")

	if v, ok := argInt(args, "limit"); ok {
		opts.Limit = v
	}
	if v, ok := argInt(args, "offset"); ok {
		opts.Offset = v
	}
	if v, ok := argBool(args, "with_deleted"); ok {
		opts.WithDeleted = v
	}

	list, err := m.store.PostList(ctx, opts)
	if err != nil {
		return "", err
	}

	items := make([]map[string]string, 0, len(list))
	for i := range list {
		post := list[i]
		items = append(items, postToMap(&post))
	}

	b, _ := json.Marshal(map[string]any{"items": items})
	return string(b), nil
}

func (m *MCP) toolPostCreate(ctx context.Context, args map[string]any) (string, error) {
	title := argString(args, "title")
	if strings.TrimSpace(title) == "" {
		return "", errors.New("title is required")
	}

	post := blogstore.NewPost()
	if id := argString(args, "id"); strings.TrimSpace(id) != "" {
		post.SetID(id)
	}

	post.SetTitle(title)

	if v := argString(args, "content"); v != "" {
		post.SetContent(v)
	}
	if v := argString(args, "summary"); v != "" {
		post.SetSummary(v)
	}
	if v := argString(args, "status"); v != "" {
		post.SetStatus(v)
	}
	if v := argString(args, "author_id"); v != "" {
		post.SetAuthorID(v)
	}
	if v := argString(args, "canonical_url"); v != "" {
		post.SetCanonicalURL(v)
	}
	if v := argString(args, "image_url"); v != "" {
		post.SetImageUrl(v)
	}
	if v := argString(args, "featured"); v != "" {
		post.SetFeatured(v)
	}
	if v := argString(args, "published_at"); v != "" {
		post.SetPublishedAt(v)
	}
	if v := argString(args, "meta_description"); v != "" {
		post.SetMetaDescription(v)
	}
	if v := argString(args, "meta_keywords"); v != "" {
		post.SetMetaKeywords(v)
	}
	if v := argString(args, "meta_robots"); v != "" {
		post.SetMetaRobots(v)
	}
	if v := argString(args, "memo"); v != "" {
		post.SetMemo(v)
	}

	if err := m.store.PostCreate(ctx, post); err != nil {
		return "", err
	}

	b, _ := json.Marshal(map[string]any{"id": post.ID(), "title": post.Title()})
	return string(b), nil
}

func (m *MCP) toolPostGet(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("id is required")
	}

	post, err := m.store.PostFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if post == nil {
		return "", errors.New("post not found")
	}

	b, _ := json.Marshal(postToMap(post))
	return string(b), nil
}

func (m *MCP) toolPostUpdate(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("id is required")
	}

	post, err := m.store.PostFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if post == nil {
		return "", errors.New("post not found")
	}

	updates := map[string]any{}
	if v, ok := args["updates"]; ok {
		if m2, ok2 := v.(map[string]any); ok2 {
			updates = m2
		}
	}

	// Support flat update arguments too
	for k, v := range args {
		if k == "id" || k == "updates" {
			continue
		}
		if _, exists := updates[k]; exists {
			continue
		}
		updates[k] = v
	}

	if v := argString(updates, "title"); v != "" {
		post.SetTitle(v)
	}
	if v := argString(updates, "content"); v != "" {
		post.SetContent(v)
	}
	if v := argString(updates, "summary"); v != "" {
		post.SetSummary(v)
	}
	if v := argString(updates, "status"); v != "" {
		post.SetStatus(v)
	}
	if v := argString(updates, "author_id"); v != "" {
		post.SetAuthorID(v)
	}
	if v := argString(updates, "canonical_url"); v != "" {
		post.SetCanonicalURL(v)
	}
	if v := argString(updates, "image_url"); v != "" {
		post.SetImageUrl(v)
	}
	if v := argString(updates, "featured"); v != "" {
		post.SetFeatured(v)
	}
	if v := argString(updates, "published_at"); v != "" {
		post.SetPublishedAt(v)
	}
	if v := argString(updates, "meta_description"); v != "" {
		post.SetMetaDescription(v)
	}
	if v := argString(updates, "meta_keywords"); v != "" {
		post.SetMetaKeywords(v)
	}
	if v := argString(updates, "meta_robots"); v != "" {
		post.SetMetaRobots(v)
	}
	if v := argString(updates, "memo"); v != "" {
		post.SetMemo(v)
	}

	if err := m.store.PostUpdate(ctx, post); err != nil {
		return "", err
	}

	b, _ := json.Marshal(map[string]any{"id": post.ID(), "title": post.Title()})
	return string(b), nil
}

func (m *MCP) toolPostDelete(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("id is required")
	}

	if err := m.store.PostDeleteByID(ctx, id); err != nil {
		return "", err
	}

	b, _ := json.Marshal(map[string]any{"deleted": true, "id": id})
	return string(b), nil
}
