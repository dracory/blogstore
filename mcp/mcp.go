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

func (m *MCP) handleInitialize(w http.ResponseWriter, _ context.Context, id any, params json.RawMessage) {
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

func (m *MCP) handleInitialized(w http.ResponseWriter, _ context.Context) {
	w.WriteHeader(http.StatusOK)
}

func (m *MCP) handleToolsList(w http.ResponseWriter, _ context.Context, id any) {
	tools := []map[string]any{
		{
			"name":        "blog_schema",
			"description": "Get schema information about blog entities and their field constraints",
			"inputSchema": map[string]any{"type": "object"},
		},
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
			"name":        "post_upsert",
			"description": "Create or update a blog post",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"title"},
				"properties": map[string]any{
					"id":               map[string]any{"type": "string"},
					"title":            map[string]any{"type": "string"},
					"content":          map[string]any{"type": "string", "description": "Post content"},
					"content_type":     map[string]any{"type": "string", "enum": []string{"markdown", "html", "plain_text"}, "default": "plain_text", "description": "Content format type for proper rendering"},
					"summary":          map[string]any{"type": "string"},
					"status":           map[string]any{"type": "string", "enum": []string{"draft", "published", "unpublished", "trash"}},
					"author_id":        map[string]any{"type": "string"},
					"canonical_url":    map[string]any{"type": "string"},
					"image_url":        map[string]any{"type": "string"},
					"featured":         map[string]any{"type": "string", "enum": []string{"yes", "no"}, "description": "Whether the post is featured (use 'yes' or 'no')"},
					"published_at":     map[string]any{"type": "string"},
					"meta_description": map[string]any{"type": "string"},
					"meta_keywords":    map[string]any{"type": "string"},
					"meta_robots":      map[string]any{"type": "string"},
					"memo":             map[string]any{"type": "string"},
				},
			},
		},
		{
			"name":        "post_versions",
			"description": "Get version history for a blog post",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"id"},
				"properties": map[string]any{
					"id":         map[string]any{"type": "string", "description": "Post ID"},
					"limit":      map[string]any{"type": "integer", "description": "Maximum number of versions to return"},
					"order_by":   map[string]any{"type": "string", "description": "Field to order by (default: created_at)"},
					"sort_order": map[string]any{"type": "string", "enum": []string{"asc", "desc"}, "description": "Sort order (default: desc)"},
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
	case "blog_schema":
		return m.toolBlogSchema(ctx, args)
	case "post_list":
		return m.toolPostList(ctx, args)
	case "post_get":
		return m.toolPostGet(ctx, args)
	case "post_upsert":
		return m.toolPostUpsert(ctx, args)
	case "post_versions":
		return m.toolPostVersions(ctx, args)
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

// contentTypeToEditor converts content_type to editor field
func contentTypeToEditor(contentType string) string {
	switch contentType {
	case blogstore.POST_CONTENT_TYPE_MARKDOWN:
		return blogstore.POST_EDITOR_MARKDOWN
	case blogstore.POST_CONTENT_TYPE_HTML:
		return blogstore.POST_EDITOR_HTMLAREA
	case blogstore.POST_CONTENT_TYPE_PLAIN_TEXT:
		return blogstore.POST_EDITOR_TEXTAREA
	default:
		return blogstore.POST_EDITOR_TEXTAREA
	}
}

func (m *MCP) toolBlogSchema(_ context.Context, _ map[string]any) (string, error) {
	schema := map[string]any{
		"entities": map[string]any{
			"post": map[string]any{
				"fields": []map[string]any{
					{"name": "id", "type": "string", "description": "Unique identifier for the post"},
					{"name": "title", "type": "string", "description": "Post title"},
					{"name": "content", "type": "string", "description": "Post content"},
					{"name": "content_type", "type": "string", "enum": []string{"markdown", "html", "plain_text"}, "description": "Content format type for proper rendering"},
					{"name": "summary", "type": "string", "description": "Brief summary of the post"},
					{"name": "status", "type": "string", "enum": []string{"draft", "published", "unpublished", "trash"}, "description": "Publication status"},
					{"name": "author_id", "type": "string", "description": "ID of the post author"},
					{"name": "canonical_url", "type": "string", "description": "Canonical URL for SEO"},
					{"name": "image_url", "type": "string", "description": "URL to featured image"},
					{"name": "featured", "type": "string", "enum": []string{"yes", "no"}, "description": "Whether the post is featured (use 'yes' or 'no')"},
					{"name": "published_at", "type": "string", "description": "Publication timestamp"},
					{"name": "meta_description", "type": "string", "description": "SEO meta description"},
					{"name": "meta_keywords", "type": "string", "description": "SEO meta keywords"},
					{"name": "meta_robots", "type": "string", "description": "SEO meta robots tag"},
					{"name": "memo", "type": "string", "description": "Internal notes"},
					{"name": "created_at", "type": "string", "description": "Creation timestamp"},
					{"name": "updated_at", "type": "string", "description": "Last update timestamp"},
					{"name": "soft_deleted_at", "type": "string", "description": "Soft deletion timestamp"},
				},
				"field_constraints": map[string]any{
					"featured": map[string]any{
						"allowed_values": []string{"yes", "no"},
						"default":        "no",
						"description":    "Must be exactly 'yes' or 'no' (not boolean true/false)",
					},
					"status": map[string]any{
						"allowed_values": []string{"draft", "published", "unpublished", "trash"},
						"default":        "draft",
						"description":    "Publication status of the post",
					},
					"content_type": map[string]any{
						"allowed_values": []string{"markdown", "html", "plain_text"},
						"default":        "plain_text",
						"description":    "Specifies how content should be rendered. Use 'markdown' for Markdown content.",
					},
					"content": map[string]any{
						"description": "Post content. The rendering is determined by the content_type field.",
					},
				},
			},
		},
		"tools": map[string]any{
			"post_list": map[string]any{
				"description": "List blog posts with filtering options",
				"arguments": map[string]any{
					"limit":        map[string]any{"type": "integer", "description": "Maximum number of posts to return"},
					"offset":       map[string]any{"type": "integer", "description": "Number of posts to skip"},
					"status":       map[string]any{"type": "string", "description": "Filter by status (draft, published, etc.)"},
					"search":       map[string]any{"type": "string", "description": "Search term for title/content"},
					"with_deleted": map[string]any{"type": "boolean", "description": "Include deleted posts"},
				},
			},
			"post_upsert": map[string]any{
				"description":        "Create or update a blog post (single operation for both create and update)",
				"required_arguments": []string{"title"},
				"arguments": map[string]any{
					"id":           map[string]any{"type": "string", "description": "Post ID (required for updates, optional for creates)"},
					"title":        map[string]any{"type": "string", "required": true, "description": "Post title"},
					"content":      map[string]any{"type": "string", "description": "Post content"},
					"content_type": map[string]any{"type": "string", "enum": []string{"markdown", "html", "plain_text"}, "default": "plain_text", "description": "Content format type for proper rendering"},
					"featured":     map[string]any{"type": "string", "enum": []string{"yes", "no"}, "default": "no", "description": "Use 'yes' or 'no' only"},
					"status":       map[string]any{"type": "string", "enum": []string{"draft", "published", "unpublished", "trash"}, "default": "draft"},
				},
			},
			"post_versions": map[string]any{
				"description":        "Get version history for a blog post (requires versioning to be enabled)",
				"required_arguments": []string{"id"},
				"arguments": map[string]any{
					"id":         map[string]any{"type": "string", "required": true, "description": "Post ID"},
					"limit":      map[string]any{"type": "integer", "description": "Maximum number of versions to return"},
					"order_by":   map[string]any{"type": "string", "description": "Field to order by (default: created_at)"},
					"sort_order": map[string]any{"type": "string", "enum": []string{"asc", "desc"}, "description": "Sort order (default: desc)"},
				},
			},
		},
		"usage_notes": []string{
			"The 'featured' field requires string values 'yes' or 'no', not boolean true/false",
			"Content supports Markdown format - use # for headers, * for emphasis, etc.",
			"Use 'published' status to make posts publicly visible",
			"Technical posts should have featured='yes' and include meta keywords",
			"Set content_type='markdown' for markdown content to enable proper rendering",
			"Use 'post_upsert' for simplified create/update operations - single method handles both cases",
			"Post updates automatically create version entries when versioning is enabled",
			"Use 'post_versions' to view and revert to previous versions of a post",
		},
	}

	result, err := json.Marshal(schema)
	if err != nil {
		return "", err
	}
	return string(result), nil
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

func (m *MCP) toolPostVersions(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("id is required")
	}

	// Check if versioning is enabled
	if !m.store.VersioningEnabled() {
		return "", errors.New("versioning is not enabled")
	}

	// Build version query
	query := blogstore.NewVersioningQuery().
		SetEntityType(blogstore.VERSIONING_TYPE_POST).
		SetEntityID(id)

	// Set optional parameters
	if orderBy := argString(args, "order_by"); orderBy != "" {
		query = query.SetOrderBy(orderBy)
	} else {
		query = query.SetOrderBy("created_at")
	}

	if sortOrder := argString(args, "sort_order"); sortOrder != "" {
		query = query.SetSortOrder(sortOrder)
	} else {
		query = query.SetSortOrder("desc")
	}

	if limit, ok := argInt(args, "limit"); ok {
		query = query.SetLimit(limit)
	}

	// Get versions
	versions, err := m.store.VersioningList(ctx, query)
	if err != nil {
		return "", err
	}

	// Convert versions to serializable format
	versionItems := make([]map[string]any, 0, len(versions))
	for _, version := range versions {
		item := map[string]any{
			"id":          version.ID(),
			"entity_id":   version.EntityID(),
			"entity_type": version.EntityType(),
			"content":     version.Content(),
			"created_at":  version.CreatedAt(),
		}
		versionItems = append(versionItems, item)
	}

	b, _ := json.Marshal(map[string]any{
		"versions": versionItems,
		"total":    len(versionItems),
	})
	return string(b), nil
}

func (m *MCP) toolPostUpsert(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	var post *blogstore.Post
	var err error
	isUpdate := false

	// Try to find existing post if ID is provided
	if strings.TrimSpace(id) != "" {
		post, err = m.store.PostFindByID(ctx, id)
		if err != nil {
			return "", err
		}
		if post != nil {
			isUpdate = true
		}
	}

	// Create new post if not found or no ID provided
	if post == nil {
		title := argString(args, "title")
		if strings.TrimSpace(title) == "" {
			return "", errors.New("title is required for new posts")
		}

		post = blogstore.NewPost()
		if strings.TrimSpace(id) != "" {
			post.SetID(id)
		}
		post.SetTitle(title)
	}

	// Set/update all fields (only if provided)
	if v := argString(args, "title"); v != "" {
		post.SetTitle(v)
	}
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
		if v != "yes" && v != "no" {
			return "", errors.New("featured field must be 'yes' or 'no', not boolean true/false")
		}
		post.SetFeatured(v)
	}

	// Set editor based on content_type
	contentType := argString(args, "content_type")
	if contentType == "" {
		// If updating existing post and no content_type provided, keep current
		if post.ContentType() == "" {
			contentType = blogstore.POST_CONTENT_TYPE_PLAIN_TEXT
		} else {
			contentType = post.ContentType()
		}
	}

	// Store content_type using the new method
	post.SetContentType(contentType)

	// Set editor based on content_type for rendering
	editor := contentTypeToEditor(contentType)
	post.SetEditor(editor)

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

	// Create or update based on whether we found an existing post
	if isUpdate {
		// Update existing post
		if err := m.store.PostUpdate(ctx, post); err != nil {
			return "", err
		}
	} else {
		// Create new post
		if err := m.store.PostCreate(ctx, post); err != nil {
			return "", err
		}
	}

	b, _ := json.Marshal(map[string]any{
		"id":     post.ID(),
		"title":  post.Title(),
		"action": "upserted",
	})
	return string(b), nil
}
