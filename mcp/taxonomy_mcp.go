package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/dracory/blogstore"
)

// ============================ TAXONOMY TOOLS ============================

func (m *MCP) taxonomyTools() []map[string]any {
	return []map[string]any{
		{
			"name":        "taxonomy_list",
			"description": "List available taxonomy types (category, tag, etc.)",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit":      map[string]any{"type": "integer"},
					"offset":     map[string]any{"type": "integer"},
					"search":     map[string]any{"type": "string"},
					"order_by":   map[string]any{"type": "string"},
					"sort_order": map[string]any{"type": "string", "enum": []string{"asc", "desc"}},
				},
			},
		},
		{
			"name":        "taxonomy_create",
			"description": "Create a new taxonomy type",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"name", "slug"},
				"properties": map[string]any{
					"name":        map[string]any{"type": "string", "description": "Display name (e.g., 'Categories')"},
					"slug":        map[string]any{"type": "string", "description": "URL slug (e.g., 'category')"},
					"description": map[string]any{"type": "string"},
				},
			},
		},
		{
			"name":        "term_list",
			"description": "List terms within a taxonomy",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"taxonomy":  map[string]any{"type": "string", "description": "Taxonomy slug (category, tag)"},
					"parent_id": map[string]any{"type": "string", "description": "Filter by parent term"},
					"search":    map[string]any{"type": "string"},
					"limit":     map[string]any{"type": "integer"},
					"offset":    map[string]any{"type": "integer"},
				},
			},
		},
		{
			"name":        "term_create",
			"description": "Create a term",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"taxonomy", "name"},
				"properties": map[string]any{
					"taxonomy":    map[string]any{"type": "string", "description": "Taxonomy slug"},
					"name":        map[string]any{"type": "string"},
					"slug":        map[string]any{"type": "string"},
					"parent_id":   map[string]any{"type": "string", "description": "For hierarchical taxonomies"},
					"description": map[string]any{"type": "string"},
				},
			},
		},
		{
			"name":        "post_set_terms",
			"description": "Set terms for a post (replaces existing)",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"post_id", "taxonomy"},
				"properties": map[string]any{
					"post_id":  map[string]any{"type": "string"},
					"taxonomy": map[string]any{"type": "string", "description": "Taxonomy slug (category, tag)"},
					"terms":    map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Term IDs"},
				},
			},
		},
		{
			"name":        "post_add_term",
			"description": "Add a single term to a post",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"post_id", "taxonomy", "term"},
				"properties": map[string]any{
					"post_id":  map[string]any{"type": "string"},
					"taxonomy": map[string]any{"type": "string"},
					"term":     map[string]any{"type": "string", "description": "Term ID"},
				},
			},
		},
		{
			"name":        "post_get_terms",
			"description": "Get terms assigned to a post",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"post_id"},
				"properties": map[string]any{
					"post_id":  map[string]any{"type": "string"},
					"taxonomy": map[string]any{"type": "string", "description": "Filter by taxonomy slug (optional)"},
				},
			},
		},
	}
}

// taxonomyToolDispatch routes taxonomy tool calls to their handlers
func (m *MCP) taxonomyToolDispatch(ctx context.Context, toolName string, args map[string]any) (string, error) {
	switch toolName {
	case "taxonomy_list":
		return m.toolTaxonomyList(ctx, args)
	case "taxonomy_create":
		return m.toolTaxonomyCreate(ctx, args)
	case "term_list":
		return m.toolTermList(ctx, args)
	case "term_create":
		return m.toolTermCreate(ctx, args)
	case "post_set_terms":
		return m.toolPostSetTerms(ctx, args)
	case "post_add_term":
		return m.toolPostAddTerm(ctx, args)
	case "post_get_terms":
		return m.toolPostGetTerms(ctx, args)
	default:
		return "", errors.New("unknown taxonomy tool")
	}
}

// toolTaxonomyList lists taxonomies
func (m *MCP) toolTaxonomyList(ctx context.Context, args map[string]any) (string, error) {
	opts := blogstore.TaxonomyQueryOptions{}
	opts.Search = argString(args, "search")
	opts.OrderBy = argString(args, "order_by")
	opts.SortOrder = argString(args, "sort_order")

	if v, ok := argInt(args, "limit"); ok {
		opts.Limit = v
	}
	if v, ok := argInt(args, "offset"); ok {
		opts.Offset = v
	}

	list, err := m.store.TaxonomyList(ctx, opts)
	if err != nil {
		return "", err
	}

	items := make([]map[string]string, 0, len(list))
	for _, t := range list {
		items = append(items, map[string]string{
			"id":          t.GetID(),
			"name":        t.GetName(),
			"slug":        t.GetSlug(),
			"description": t.GetDescription(),
			"created_at":  t.GetCreatedAt(),
		})
	}

	b, _ := json.Marshal(map[string]any{"items": items})
	return string(b), nil
}

// toolTaxonomyCreate creates a new taxonomy
func (m *MCP) toolTaxonomyCreate(ctx context.Context, args map[string]any) (string, error) {
	name := argString(args, "name")
	slug := argString(args, "slug")
	description := argString(args, "description")

	if strings.TrimSpace(name) == "" || strings.TrimSpace(slug) == "" {
		return "", errors.New("name and slug are required")
	}

	taxonomy := blogstore.NewTaxonomy()
	taxonomy.SetName(name).
		SetSlug(slug).
		SetDescription(description)

	if err := m.store.TaxonomyCreate(ctx, taxonomy); err != nil {
		return "", err
	}

	b, _ := json.Marshal(map[string]any{
		"id":     taxonomy.GetID(),
		"name":   taxonomy.GetName(),
		"slug":   taxonomy.GetSlug(),
		"action": "created",
	})
	return string(b), nil
}

// toolTermList lists terms
func (m *MCP) toolTermList(ctx context.Context, args map[string]any) (string, error) {
	opts := blogstore.TermQueryOptions{}
	opts.Search = argString(args, "search")
	opts.ParentID = argString(args, "parent_id")

	if v, ok := argInt(args, "limit"); ok {
		opts.Limit = v
	}
	if v, ok := argInt(args, "offset"); ok {
		opts.Offset = v
	}

	// If taxonomy slug provided, look up the taxonomy ID
	taxonomySlug := argString(args, "taxonomy")
	if taxonomySlug != "" {
		taxonomy, err := m.store.TaxonomyFindBySlug(ctx, taxonomySlug)
		if err != nil {
			return "", err
		}
		if taxonomy == nil {
			return "", errors.New("taxonomy not found: " + taxonomySlug)
		}
		opts.TaxonomyID = taxonomy.GetID()
	}

	list, err := m.store.TermList(ctx, opts)
	if err != nil {
		return "", err
	}

	items := make([]map[string]any, 0, len(list))
	for _, t := range list {
		items = append(items, map[string]any{
			"id":          t.GetID(),
			"taxonomy_id": t.GetTaxonomyID(),
			"parent_id":   t.GetParentID(),
			"name":        t.GetName(),
			"slug":        t.GetSlug(),
			"description": t.GetDescription(),
			"count":       t.GetCount(),
			"created_at":  t.GetCreatedAt(),
		})
	}

	b, _ := json.Marshal(map[string]any{"items": items})
	return string(b), nil
}

// toolTermCreate creates a new term
func (m *MCP) toolTermCreate(ctx context.Context, args map[string]any) (string, error) {
	taxonomySlug := argString(args, "taxonomy")
	name := argString(args, "name")
	slug := argString(args, "slug")
	parentID := argString(args, "parent_id")
	description := argString(args, "description")

	if strings.TrimSpace(taxonomySlug) == "" || strings.TrimSpace(name) == "" {
		return "", errors.New("taxonomy and name are required")
	}

	// Find taxonomy
	taxonomy, err := m.store.TaxonomyFindBySlug(ctx, taxonomySlug)
	if err != nil {
		return "", err
	}
	if taxonomy == nil {
		return "", errors.New("taxonomy not found: " + taxonomySlug)
	}

	term := blogstore.NewTerm()
	term.SetTaxonomyID(taxonomy.GetID()).
		SetName(name).
		SetSlug(slug).
		SetParentID(parentID).
		SetDescription(description)

	if err := m.store.TermCreate(ctx, term); err != nil {
		return "", err
	}

	b, _ := json.Marshal(map[string]any{
		"id":          term.GetID(),
		"taxonomy_id": term.GetTaxonomyID(),
		"name":        term.GetName(),
		"slug":        term.GetSlug(),
		"action":      "created",
	})
	return string(b), nil
}

// toolPostSetTerms sets terms for a post (replaces existing)
func (m *MCP) toolPostSetTerms(ctx context.Context, args map[string]any) (string, error) {
	postID := argString(args, "post_id")
	taxonomySlug := argString(args, "taxonomy")

	if strings.TrimSpace(postID) == "" || strings.TrimSpace(taxonomySlug) == "" {
		return "", errors.New("post_id and taxonomy are required")
	}

	// Verify post exists
	post, err := m.store.PostFindByID(ctx, postID)
	if err != nil {
		return "", err
	}
	if post == nil {
		return "", errors.New("post not found: " + postID)
	}

	// Get term IDs from args
	var termIDs []string
	if terms, ok := args["terms"].([]any); ok {
		for _, t := range terms {
			if s, ok := t.(string); ok && s != "" {
				termIDs = append(termIDs, s)
			}
		}
	}

	// Set terms via store
	if err := m.store.PostSetTerms(ctx, postID, taxonomySlug, termIDs); err != nil {
		return "", err
	}

	// Update post metadata
	post.SetTermIDs(taxonomySlug, termIDs)
	if err := m.store.PostUpdate(ctx, post); err != nil {
		return "", err
	}

	b, _ := json.Marshal(map[string]any{
		"post_id":  postID,
		"taxonomy": taxonomySlug,
		"terms":    termIDs,
		"action":   "set_terms",
	})
	return string(b), nil
}

// toolPostAddTerm adds a single term to a post
func (m *MCP) toolPostAddTerm(ctx context.Context, args map[string]any) (string, error) {
	postID := argString(args, "post_id")
	taxonomySlug := argString(args, "taxonomy")
	termID := argString(args, "term")

	if strings.TrimSpace(postID) == "" || strings.TrimSpace(taxonomySlug) == "" || strings.TrimSpace(termID) == "" {
		return "", errors.New("post_id, taxonomy, and term are required")
	}

	// Verify post exists
	post, err := m.store.PostFindByID(ctx, postID)
	if err != nil {
		return "", err
	}
	if post == nil {
		return "", errors.New("post not found: " + postID)
	}

	// Add term via store
	if err := m.store.PostInsertTermAt(ctx, postID, termID, 0); err != nil {
		return "", err
	}

	// Update post metadata
	existingIDs := post.TermIDs(taxonomySlug)
	found := false
	for _, id := range existingIDs {
		if id == termID {
			found = true
			break
		}
	}
	if !found {
		newIDs := append(existingIDs, termID)
		post.SetTermIDs(taxonomySlug, newIDs)
		if err := m.store.PostUpdate(ctx, post); err != nil {
			return "", err
		}
	}

	b, _ := json.Marshal(map[string]any{
		"post_id":  postID,
		"taxonomy": taxonomySlug,
		"term_id":  termID,
		"action":   "added",
	})
	return string(b), nil
}

// toolPostGetTerms gets terms assigned to a post
func (m *MCP) toolPostGetTerms(ctx context.Context, args map[string]any) (string, error) {
	postID := argString(args, "post_id")
	taxonomySlug := argString(args, "taxonomy")

	if strings.TrimSpace(postID) == "" {
		return "", errors.New("post_id is required")
	}

	// Get terms via store
	terms, err := m.store.TermListByPostID(ctx, postID, taxonomySlug)
	if err != nil {
		return "", err
	}

	items := make([]map[string]any, 0, len(terms))
	for _, t := range terms {
		items = append(items, map[string]any{
			"id":          t.GetID(),
			"taxonomy_id": t.GetTaxonomyID(),
			"parent_id":   t.GetParentID(),
			"name":        t.GetName(),
			"slug":        t.GetSlug(),
			"count":       t.GetCount(),
		})
	}

	b, _ := json.Marshal(map[string]any{
		"post_id":  postID,
		"taxonomy": taxonomySlug,
		"terms":    items,
	})
	return string(b), nil
}
