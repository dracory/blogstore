package blogstore

import (
	"context"
	"errors"
	"time"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

// ============================ TAXONOMY STORE METHODS ============================

// TaxonomyCount returns the total number of taxonomies matching the given query options.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyCount(ctx context.Context, options TaxonomyQueryOptions) (int64, error) {
	if ctx == nil {
		return 0, errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return -1, errors.New("taxonomy is not enabled")
	}

	q := store.buildTaxonomyQuery(options)

	var count int64
	err := q.Table(store.taxonomyTableName).Count(&count)
	return count, err
}

// buildTaxonomyQuery builds a neat query from the taxonomy query options.
func (store *storeImplementation) buildTaxonomyQuery(options TaxonomyQueryOptions) contractsorm.Query {
	q := store.db.Query()

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if options.Slug != "" {
		q = q.Where(COLUMN_SLUG+" = ?", options.Slug)
	}

	if options.Limit > 0 {
		q = q.Limit(options.Limit)
	}

	if options.Offset > 0 {
		q = q.Offset(options.Offset)
	}

	return q
}

// TaxonomyCreate inserts a new taxonomy into the database.
// Sets the created_at and updated_at timestamps automatically.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyCreate(ctx context.Context, taxonomy TaxonomyInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if taxonomy == nil {
		return errors.New("taxonomy is nil")
	}

	if taxonomy.GetID() == "" {
		taxonomy.SetID(GenerateShortID())
	}

	if taxonomy.GetCreatedAt() == "" {
		taxonomy.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if taxonomy.GetUpdatedAt() == "" {
		taxonomy.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	data := taxonomy.GetData()
	delete(data, COLUMN_ID)
	delete(data, COLUMN_CREATED_AT)
	delete(data, COLUMN_UPDATED_AT)

	db, err := store.db.DB()
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "INSERT INTO "+store.taxonomyTableName+" (id, name, slug, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		taxonomy.GetID(),
		taxonomy.GetName(),
		taxonomy.GetSlug(),
		taxonomy.GetDescription(),
		taxonomy.GetCreatedAtCarbon().StdTime(),
		taxonomy.GetUpdatedAtCarbon().StdTime(),
	)

	return err
}

// TaxonomyDelete permanently removes a taxonomy from the database.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyDelete(ctx context.Context, taxonomy TaxonomyInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if taxonomy == nil {
		return errors.New("taxonomy is nil")
	}

	return store.TaxonomyDeleteByID(ctx, taxonomy.GetID())
}

// TaxonomyDeleteByID permanently removes a taxonomy by its ID.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyDeleteByID(ctx context.Context, id string) error {
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if id == "" {
		return errors.New("taxonomy id is empty")
	}

	_, err := store.db.Query().
		Table(store.taxonomyTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

// TaxonomyFindByID retrieves a taxonomy by its ID.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyFindByID(ctx context.Context, id string) (TaxonomyInterface, error) {
	if !store.taxonomyEnabled {
		return nil, errors.New("taxonomy is not enabled")
	}
	if id == "" {
		return nil, errors.New("taxonomy id is empty")
	}

	list, err := store.TaxonomyList(ctx, TaxonomyQueryOptions{
		ID:    id,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// TaxonomyFindBySlug retrieves a taxonomy by its slug.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyFindBySlug(ctx context.Context, slug string) (TaxonomyInterface, error) {
	if !store.taxonomyEnabled {
		return nil, errors.New("taxonomy is not enabled")
	}
	if slug == "" {
		return nil, errors.New("taxonomy slug is empty")
	}

	list, err := store.TaxonomyList(ctx, TaxonomyQueryOptions{
		Slug:  slug,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// TaxonomyList retrieves a list of taxonomies matching the given query options.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyInterface, error) {
	if !store.taxonomyEnabled {
		return []TaxonomyInterface{}, errors.New("taxonomy is not enabled")
	}
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	type taxonomyRow struct {
		ID          string    `db:"id"`
		Name        string    `db:"name"`
		Slug        string    `db:"slug"`
		Description string    `db:"description"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}

	q := store.buildTaxonomyQuery(options)

	var rows []taxonomyRow
	if err := q.Table(store.taxonomyTableName).Get(&rows); err != nil {
		return []TaxonomyInterface{}, err
	}

	list := make([]TaxonomyInterface, 0, len(rows))
	for _, r := range rows {
		t := NewTaxonomy()
		t.SetID(r.ID)
		t.SetName(r.Name)
		t.SetSlug(r.Slug)
		t.SetDescription(r.Description)
		if taxImpl, ok := t.(*taxonomyImplementation); ok {
			taxImpl.CreatedAtField.CreatedAt = r.CreatedAt
			taxImpl.UpdatedAtField.UpdatedAt = r.UpdatedAt
		}
		list = append(list, t)
	}

	return list, nil
}

// TaxonomyUpdate updates an existing taxonomy in the database.
// Only changed fields are updated. Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyUpdate(ctx context.Context, taxonomy TaxonomyInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if taxonomy == nil {
		return errors.New("taxonomy is nil")
	}

	taxonomy.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	_, err := store.db.Query().
		Table(store.taxonomyTableName).
		Where(COLUMN_ID+" = ?", taxonomy.GetID()).
		Update(map[string]interface{}{
			COLUMN_NAME:        taxonomy.GetName(),
			COLUMN_SLUG:        taxonomy.GetSlug(),
			COLUMN_DESCRIPTION: taxonomy.GetDescription(),
			COLUMN_UPDATED_AT:  taxonomy.GetUpdatedAtCarbon().StdTime(),
		})

	return err
}

// ============================ TERM STORE METHODS ============================

// TermCount returns the total number of terms matching the given query options.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermCount(ctx context.Context, options TermQueryOptions) (int64, error) {
	if ctx == nil {
		return 0, errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return -1, errors.New("taxonomy is not enabled")
	}

	q := store.buildTermQuery(options)

	var count int64
	err := q.Table(store.termTableName).Count(&count)
	return count, err
}

// buildTermQuery builds a neat query from the term query options.
func (store *storeImplementation) buildTermQuery(options TermQueryOptions) contractsorm.Query {
	q := store.db.Query()

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if len(options.IDIn) > 0 {
		// Build IN clause manually for neat compatibility
		inClause := COLUMN_ID + " IN ("
		placeholders := make([]interface{}, 0, len(options.IDIn))
		for i, id := range options.IDIn {
			if i > 0 {
				inClause += ", "
			}
			inClause += "?"
			placeholders = append(placeholders, id)
		}
		inClause += ")"
		q = q.Where(inClause, placeholders...)
	}

	if options.TaxonomyID != "" {
		q = q.Where(COLUMN_TAXONOMY_ID+" = ?", options.TaxonomyID)
	}

	if options.ParentID != "" {
		q = q.Where(COLUMN_PARENT_ID+" = ?", options.ParentID)
	}

	if options.Slug != "" {
		q = q.Where(COLUMN_SLUG+" = ?", options.Slug)
	}

	if options.Limit > 0 {
		q = q.Limit(options.Limit)
	}

	if options.Offset > 0 {
		q = q.Offset(options.Offset)
	}

	return q
}

// TermCreate inserts a new term into the database.
// Sets the created_at and updated_at timestamps automatically.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermCreate(ctx context.Context, term TermInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if term == nil {
		return errors.New("term is nil")
	}

	if term.GetID() == "" {
		term.SetID(GenerateShortID())
	}

	if term.GetCreatedAt() == "" {
		term.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if term.GetUpdatedAt() == "" {
		term.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	db, err := store.db.DB()
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "INSERT INTO "+store.termTableName+" (id, taxonomy_id, parent_id, sequence, name, slug, description, count, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		term.GetID(),
		term.GetTaxonomyID(),
		term.GetParentID(),
		term.GetSequence(),
		term.GetName(),
		term.GetSlug(),
		term.GetDescription(),
		term.GetCount(),
		term.GetCreatedAtCarbon().StdTime(),
		term.GetUpdatedAtCarbon().StdTime(),
	)

	return err
}

// TermDelete permanently removes a term from the database.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermDelete(ctx context.Context, term TermInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if term == nil {
		return errors.New("term is nil")
	}

	return store.TermDeleteByID(ctx, term.GetID())
}

// TermDeleteByID permanently removes a term by its ID.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermDeleteByID(ctx context.Context, id string) error {
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if id == "" {
		return errors.New("term id is empty")
	}

	_, err := store.db.Query().
		Table(store.termTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

// TermFindByID retrieves a term by its ID.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermFindByID(ctx context.Context, id string) (TermInterface, error) {
	if !store.taxonomyEnabled {
		return nil, errors.New("taxonomy is not enabled")
	}
	if id == "" {
		return nil, errors.New("term id is empty")
	}

	list, err := store.TermList(ctx, TermQueryOptions{
		ID:    id,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// TermFindBySlug retrieves a term by its taxonomy slug and term slug.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermFindBySlug(ctx context.Context, taxonomySlug, termSlug string) (TermInterface, error) {
	if !store.taxonomyEnabled {
		return nil, errors.New("taxonomy is not enabled")
	}
	if taxonomySlug == "" || termSlug == "" {
		return nil, errors.New("taxonomy slug and term slug are required")
	}

	taxonomy, err := store.TaxonomyFindBySlug(ctx, taxonomySlug)
	if err != nil {
		return nil, err
	}
	if taxonomy == nil {
		return nil, nil
	}

	list, err := store.TermList(ctx, TermQueryOptions{
		TaxonomyID: taxonomy.GetID(),
		Slug:       termSlug,
		Limit:      1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// TermList retrieves a list of terms matching the given query options.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermList(ctx context.Context, options TermQueryOptions) ([]TermInterface, error) {
	if !store.taxonomyEnabled {
		return []TermInterface{}, errors.New("taxonomy is not enabled")
	}
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	type termRow struct {
		ID          string    `db:"id"`
		TaxonomyID  string    `db:"taxonomy_id"`
		ParentID    string    `db:"parent_id"`
		Sequence    int       `db:"sequence"`
		Name        string    `db:"name"`
		Slug        string    `db:"slug"`
		Description string    `db:"description"`
		Count       int       `db:"count"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}

	q := store.buildTermQuery(options)

	var rows []termRow
	if err := q.Table(store.termTableName).Get(&rows); err != nil {
		return []TermInterface{}, err
	}

	list := make([]TermInterface, 0, len(rows))
	for _, r := range rows {
		t := NewTerm()
		t.SetID(r.ID)
		t.SetTaxonomyID(r.TaxonomyID)
		t.SetParentID(r.ParentID)
		t.SetSequence(r.Sequence)
		t.SetName(r.Name)
		t.SetSlug(r.Slug)
		t.SetDescription(r.Description)
		t.SetCount(r.Count)
		if termImpl, ok := t.(*termImplementation); ok {
			termImpl.CreatedAtField.CreatedAt = r.CreatedAt
			termImpl.UpdatedAtField.UpdatedAt = r.UpdatedAt
		}
		list = append(list, t)
	}

	return list, nil
}

// TermUpdate updates an existing term in the database.
// Only changed fields are updated. Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermUpdate(ctx context.Context, term TermInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if term == nil {
		return errors.New("term is nil")
	}

	term.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	_, err := store.db.Query().
		Table(store.termTableName).
		Where(COLUMN_ID+" = ?", term.GetID()).
		Update(map[string]interface{}{
			COLUMN_TAXONOMY_ID: term.GetTaxonomyID(),
			COLUMN_PARENT_ID:   term.GetParentID(),
			COLUMN_SEQUENCE:    term.GetSequence(),
			COLUMN_NAME:        term.GetName(),
			COLUMN_SLUG:        term.GetSlug(),
			COLUMN_DESCRIPTION: term.GetDescription(),
			COLUMN_COUNT:       term.GetCount(),
			COLUMN_UPDATED_AT:  term.GetUpdatedAtCarbon().StdTime(),
		})

	return err
}

// ============================ POST-TERM RELATIONSHIP METHODS ============================

// PostInsertTermAt creates a relationship between a post and a term at a specific sequence position.
// Also increments the term's count. Duplicate key errors are ignored.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) PostInsertTermAt(ctx context.Context, postID string, termID string, sequence int) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if postID == "" || termID == "" {
		return errors.New("post id and term id are required")
	}

	relationID := GenerateShortID()
	now := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)

	db, err := store.db.DB()
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "INSERT INTO "+store.termRelationTableName+" (id, post_id, term_id, sequence, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		relationID,
		postID,
		termID,
		sequence,
		now,
		now,
	)

	if err != nil {
		return err
	}

	// Increment term count
	_, err = db.ExecContext(ctx, "UPDATE "+store.termTableName+" SET count = count + 1 WHERE id = ?", termID)
	return err
}

// PostAddTerm appends a term to a post at the end of the sequence.
// Automatically calculates the next available sequence number.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) PostAddTerm(ctx context.Context, postID string, termID string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if postID == "" || termID == "" {
		return errors.New("post id and term id are required")
	}

	maxSeq, err := store.postMaxSequence(ctx, postID)
	if err != nil {
		return err
	}

	return store.PostInsertTermAt(ctx, postID, termID, maxSeq+1)
}

// postMaxSequence returns the maximum sequence number for terms associated with a post.
// Returns 0 if no terms are associated with the post.
func (store *storeImplementation) postMaxSequence(ctx context.Context, postID string) (int, error) {
	db, err := store.db.DB()
	if err != nil {
		return 0, err
	}
	var maxSeq int
	err = db.QueryRowContext(ctx, "SELECT COALESCE(MAX(sequence), 0) FROM "+store.termRelationTableName+" WHERE post_id = ?", postID).Scan(&maxSeq)
	return maxSeq, err
}

// PostMoveTermTo moves a term to a specific sequence position on a post.
// Reorders existing terms by fetching all, reordering in memory, and updating one by one.
// Returns an error if the term is not associated with the post.
func (store *storeImplementation) PostMoveTermTo(ctx context.Context, postID string, termID string, sequence int) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if postID == "" || termID == "" {
		return errors.New("post id and term id are required")
	}
	if sequence < 0 {
		return errors.New("sequence must be non-negative")
	}

	relations, err := store.postTermList(ctx, postID)
	if err != nil {
		return err
	}

	// Find the term relation
	var targetRelation TermRelationInterface
	var targetIndex int
	for i, rel := range relations {
		if rel.GetTermID() == termID {
			targetRelation = rel
			targetIndex = i
			break
		}
	}

	if targetRelation == nil {
		return errors.New("term not associated with post")
	}

	// Remove from current position
	relations = append(relations[:targetIndex], relations[targetIndex+1:]...)

	// Convert 1-based sequence to 0-based for slice operations
	zeroBasedIndex := sequence - 1
	if zeroBasedIndex < 0 {
		zeroBasedIndex = 0
	}

	// Insert at new position
	if zeroBasedIndex >= len(relations) {
		relations = append(relations, targetRelation)
	} else {
		relations = append(relations[:zeroBasedIndex], append([]TermRelationInterface{targetRelation}, relations[zeroBasedIndex:]...)...)
	}

	// Update all sequences (1-based indexing)
	for i, rel := range relations {
		if err := store.postTermUpdateSequence(ctx, postID, rel.GetTermID(), i+1); err != nil {
			return err
		}
	}

	return nil
}

// postTermList returns all term relations for a post.
func (store *storeImplementation) postTermList(ctx context.Context, postID string) ([]TermRelationInterface, error) {
	type termRelationRow struct {
		ID        string    `db:"id"`
		PostID    string    `db:"post_id"`
		TermID    string    `db:"term_id"`
		Sequence  int       `db:"sequence"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	db, err := store.db.DB()
	if err != nil {
		return []TermRelationInterface{}, err
	}
	rows, err := db.QueryContext(ctx, "SELECT id, post_id, term_id, sequence, created_at, updated_at FROM "+store.termRelationTableName+" WHERE post_id = ? ORDER BY sequence ASC", postID)
	if err != nil {
		return []TermRelationInterface{}, err
	}
	defer rows.Close()

	var rowList []termRelationRow
	for rows.Next() {
		var r termRelationRow
		if err := rows.Scan(&r.ID, &r.PostID, &r.TermID, &r.Sequence, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return []TermRelationInterface{}, err
		}
		rowList = append(rowList, r)
	}
	if err := rows.Err(); err != nil {
		return []TermRelationInterface{}, err
	}

	list := make([]TermRelationInterface, 0, len(rowList))
	for _, r := range rowList {
		tr := NewTermRelation()
		tr.SetID(r.ID)
		tr.SetPostID(r.PostID)
		tr.SetTermID(r.TermID)
		tr.SetSequence(r.Sequence)
		if trImpl, ok := tr.(*termRelationImplementation); ok {
			trImpl.CreatedAtField.CreatedAt = r.CreatedAt
			trImpl.UpdatedAtField.UpdatedAt = r.UpdatedAt
		}
		list = append(list, tr)
	}

	return list, nil
}

// postTermUpdateSequence updates the sequence of a specific term relation.
func (store *storeImplementation) postTermUpdateSequence(ctx context.Context, postID, termID string, sequence int) error {
	_, err := store.db.Query().
		Table(store.termRelationTableName).
		Where(COLUMN_POST_ID+" = ? AND "+COLUMN_TERM_ID+" = ?", postID, termID).
		Update(map[string]interface{}{
			COLUMN_SEQUENCE:   sequence,
			COLUMN_UPDATED_AT: carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		})
	return err
}

// PostRemoveTerm removes the relationship between a post and a term.
// Also decrements the term's count.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) PostRemoveTerm(ctx context.Context, postID string, termID string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if postID == "" || termID == "" {
		return errors.New("post id and term id are required")
	}

	_, err := store.db.Query().
		Table(store.termRelationTableName).
		Where(COLUMN_POST_ID+" = ? AND "+COLUMN_TERM_ID+" = ?", postID, termID).
		Delete()
	if err != nil {
		return err
	}

	// Decrement term count
	db, err := store.db.DB()
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "UPDATE "+store.termTableName+" SET count = CASE WHEN count > 0 THEN count - 1 ELSE 0 END WHERE id = ?", termID)
	return err
}

// TermListByPostID retrieves all terms associated with a specific post.
// Optionally filters by taxonomy slug.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermListByPostID(ctx context.Context, postID string, taxonomySlug string) ([]TermInterface, error) {
	if !store.taxonomyEnabled {
		return []TermInterface{}, errors.New("taxonomy is not enabled")
	}
	if postID == "" {
		return []TermInterface{}, errors.New("post id is required")
	}

	relations, err := store.postTermList(ctx, postID)
	if err != nil {
		return []TermInterface{}, err
	}

	termIDs := make([]string, 0, len(relations))
	for _, rel := range relations {
		termIDs = append(termIDs, rel.GetTermID())
	}

	if len(termIDs) == 0 {
		return []TermInterface{}, nil
	}

	terms, err := store.TermList(ctx, TermQueryOptions{
		IDIn: termIDs,
	})
	if err != nil {
		return []TermInterface{}, err
	}

	// Filter by taxonomy slug if provided
	if taxonomySlug != "" {
		taxonomy, err := store.TaxonomyFindBySlug(ctx, taxonomySlug)
		if err != nil {
			return []TermInterface{}, err
		}
		if taxonomy == nil {
			return []TermInterface{}, nil
		}

		filtered := make([]TermInterface, 0)
		for _, term := range terms {
			if term.GetTaxonomyID() == taxonomy.GetID() {
				filtered = append(filtered, term)
			}
		}
		return filtered, nil
	}

	return terms, nil
}

// PostListByTermID retrieves all posts associated with a specific term.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) PostListByTermID(ctx context.Context, termID string, options PostQueryOptions) ([]PostInterface, error) {
	if ctx == nil {
		return []PostInterface{}, errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return []PostInterface{}, errors.New("taxonomy is not enabled")
	}
	if termID == "" {
		return []PostInterface{}, errors.New("term id is required")
	}

	db, err := store.db.DB()
	if err != nil {
		return []PostInterface{}, err
	}
	rows, err := db.QueryContext(ctx, "SELECT post_id FROM "+store.termRelationTableName+" WHERE term_id = ?", termID)
	if err != nil {
		return []PostInterface{}, err
	}
	defer rows.Close()

	var postIDs []string
	for rows.Next() {
		var postID string
		if err := rows.Scan(&postID); err != nil {
			return []PostInterface{}, err
		}
		postIDs = append(postIDs, postID)
	}

	if len(postIDs) == 0 {
		return []PostInterface{}, nil
	}

	options.IDIn = postIDs
	return store.PostList(ctx, options)
}

// PostSetTerms sets the terms for a post within a specific taxonomy.
// Removes any existing terms not in the provided list and adds new ones.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) PostSetTerms(ctx context.Context, postID string, taxonomySlug string, termIDs []string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if postID == "" {
		return errors.New("post id is required")
	}

	taxonomy, err := store.TaxonomyFindBySlug(ctx, taxonomySlug)
	if err != nil {
		return err
	}
	if taxonomy == nil {
		return errors.New("taxonomy not found")
	}

	// Get current terms for this post and taxonomy
	currentTerms, err := store.TermListByPostID(ctx, postID, taxonomySlug)
	if err != nil {
		return err
	}

	// Build map of current term IDs
	currentTermIDs := make(map[string]bool)
	for _, term := range currentTerms {
		currentTermIDs[term.GetID()] = true
	}

	// Build map of new term IDs
	newTermIDs := make(map[string]bool)
	for _, termID := range termIDs {
		newTermIDs[termID] = true
	}

	// Remove terms that are no longer in the list
	for _, term := range currentTerms {
		if !newTermIDs[term.GetID()] {
			if err := store.PostRemoveTerm(ctx, postID, term.GetID()); err != nil {
				return err
			}
		}
	}

	// Add new terms
	for _, termID := range termIDs {
		if !currentTermIDs[termID] {
			if err := store.PostAddTerm(ctx, postID, termID); err != nil {
				return err
			}
		}
	}

	return nil
}

// ============================ UTILITY METHODS ============================

// TermIncrementCount increments the count of posts associated with a term.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermIncrementCount(ctx context.Context, termID string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if termID == "" {
		return errors.New("term id is required")
	}

	db, err := store.db.DB()
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "UPDATE "+store.termTableName+" SET count = count + 1 WHERE id = ?", termID)
	return err
}

// TermDecrementCount decrements the count of posts associated with a term.
// The count will not go below zero.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermDecrementCount(ctx context.Context, termID string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if termID == "" {
		return errors.New("term id is required")
	}

	db, err := store.db.DB()
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "UPDATE "+store.termTableName+" SET count = CASE WHEN count > 0 THEN count - 1 ELSE 0 END WHERE id = ?", termID)
	return err
}
