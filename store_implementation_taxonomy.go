package blogstore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

// ============================ TAXONOMY STORE METHODS ============================

// TaxonomyCount returns the total number of taxonomies matching the given query options.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyCount(ctx context.Context, options TaxonomyQueryOptions) (int64, error) {
	if !store.taxonomyEnabled {
		return -1, errors.New("taxonomy is not enabled")
	}
	options.CountOnly = true
	q := store.taxonomyQuery(options)

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	mapped, err := database.SelectToMapString(
		database.NewQueryableContext(ctx, store.db),
		sqlStr,
		params...,
	)

	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]
	i, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		return -1, err
	}

	return i, nil
}

// TaxonomyCreate inserts a new taxonomy into the database.
// Sets the created_at and updated_at timestamps automatically.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyCreate(ctx context.Context, taxonomy TaxonomyInterface) error {
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	taxonomy.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	data := taxonomy.GetData()

	sqlStr, sqlParams, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.taxonomyTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	taxonomy.MarkAsNotDirty()
	return nil
}

// TaxonomyDelete permanently removes a taxonomy from the database.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyDelete(ctx context.Context, taxonomy TaxonomyInterface) error {
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

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.taxonomyTableName).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		Prepared(true).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

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
	q := store.taxonomyQuery(options)

	sqlStr, sqlParams, errSql := q.Select().
		Prepared(true).
		ToSQL()

	if errSql != nil {
		log.Println(errSql)
		return []TaxonomyInterface{}, errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(
		database.NewQueryableContext(ctx, store.db),
		sqlStr,
		sqlParams...,
	)
	if err != nil {
		return []TaxonomyInterface{}, err
	}

	list := []TaxonomyInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewTaxonomyFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

// TaxonomyUpdate updates an existing taxonomy in the database.
// Only changed fields are updated. Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TaxonomyUpdate(ctx context.Context, taxonomy TaxonomyInterface) error {
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if taxonomy == nil {
		return errors.New("taxonomy is nil")
	}

	dataChanged := taxonomy.GetDataChanged()

	delete(dataChanged, "id")

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.taxonomyTableName).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(taxonomy.GetID())).
		Prepared(true).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	taxonomy.MarkAsNotDirty()

	return err
}

// taxonomyQuery builds a goqu SelectDataset for querying taxonomies based on options.
func (store *storeImplementation) taxonomyQuery(options TaxonomyQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).
		From(store.taxonomyTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Slug != "" {
		q = q.Where(goqu.C(COLUMN_SLUG).Eq(options.Slug))
	}

	if options.Search != "" {
		var searchExpr goqu.Expression
		switch store.dbDriverName {
		case "sqlite3", "sqlite":
			// SQLite: use LOWER() for case-insensitive matching
			searchPattern := "%" + options.Search + "%"
			searchExpr = goqu.Or(
				goqu.L("LOWER(?)", goqu.C(COLUMN_NAME)).Like(goqu.L("LOWER(?)", searchPattern)),
				goqu.L("LOWER(?)", goqu.C(COLUMN_DESCRIPTION)).Like(goqu.L("LOWER(?)", searchPattern)),
			)
		default:
			// PostgreSQL, MySQL: use ILike
			searchExpr = goqu.Or(
				goqu.C(COLUMN_NAME).ILike("%"+options.Search+"%"),
				goqu.C(COLUMN_DESCRIPTION).ILike("%"+options.Search+"%"),
			)
		}
		q = q.Where(searchExpr)
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}

		sortOrder := "asc"
		if options.SortOrder != "" {
			sortOrder = options.SortOrder
		}

		orderBy := COLUMN_NAME
		if options.OrderBy != "" {
			orderBy = options.OrderBy
		}

		if sortOrder == sb.ASC {
			q = q.Order(goqu.I(orderBy).Asc())
		} else {
			q = q.Order(goqu.I(orderBy).Desc())
		}
	}

	return q
}

// ============================ TERM STORE METHODS ============================

// TermCount returns the total number of terms matching the given query options.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermCount(ctx context.Context, options TermQueryOptions) (int64, error) {
	if !store.taxonomyEnabled {
		return -1, errors.New("taxonomy is not enabled")
	}
	options.CountOnly = true
	q := store.termQuery(options)

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	mapped, err := database.SelectToMapString(
		database.NewQueryableContext(ctx, store.db),
		sqlStr,
		params...,
	)

	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]
	i, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		return -1, err
	}

	return i, nil
}

// TermCreate inserts a new term into the database.
// Sets the created_at and updated_at timestamps automatically.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermCreate(ctx context.Context, term TermInterface) error {
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	term.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	data := term.GetData()

	sqlStr, sqlParams, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.termTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	term.MarkAsNotDirty()
	return nil
}

// TermDelete permanently removes a term from the database.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermDelete(ctx context.Context, term TermInterface) error {
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

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.termTableName).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		Prepared(true).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

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
		Search:     termSlug,
		Limit:      1,
	})

	if err != nil {
		return nil, err
	}

	// Filter by exact slug match
	for _, term := range list {
		if term.GetSlug() == termSlug {
			return term, nil
		}
	}

	return nil, nil
}

// TermList retrieves a list of terms matching the given query options.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermList(ctx context.Context, options TermQueryOptions) ([]TermInterface, error) {
	if !store.taxonomyEnabled {
		return []TermInterface{}, errors.New("taxonomy is not enabled")
	}
	q := store.termQuery(options)

	sqlStr, sqlParams, errSql := q.Select().
		Prepared(true).
		ToSQL()

	if errSql != nil {
		log.Println(errSql)
		return []TermInterface{}, errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(
		database.NewQueryableContext(ctx, store.db),
		sqlStr,
		sqlParams...,
	)
	if err != nil {
		return []TermInterface{}, err
	}

	list := []TermInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewTermFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

// TermUpdate updates an existing term in the database.
// Only changed fields are updated. Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermUpdate(ctx context.Context, term TermInterface) error {
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if term == nil {
		return errors.New("term is nil")
	}

	dataChanged := term.GetDataChanged()

	delete(dataChanged, "id")

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.termTableName).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(term.GetID())).
		Prepared(true).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	term.MarkAsNotDirty()

	return err
}

// termQuery builds a goqu SelectDataset for querying terms based on options.
func (store *storeImplementation) termQuery(options TermQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).
		From(store.termTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if len(options.IDIn) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn))
	}

	if options.TaxonomyID != "" {
		q = q.Where(goqu.C(COLUMN_TAXONOMY_ID).Eq(options.TaxonomyID))
	}

	if options.ParentID != "" {
		q = q.Where(goqu.C(COLUMN_PARENT_ID).Eq(options.ParentID))
	}

	if options.Search != "" {
		var searchExpr goqu.Expression
		switch store.dbDriverName {
		case "sqlite3", "sqlite":
			// SQLite: use LOWER() for case-insensitive matching
			searchPattern := "%" + options.Search + "%"
			searchExpr = goqu.Or(
				goqu.L("LOWER(?)", goqu.C(COLUMN_NAME)).Like(goqu.L("LOWER(?)", searchPattern)),
				goqu.L("LOWER(?)", goqu.C(COLUMN_DESCRIPTION)).Like(goqu.L("LOWER(?)", searchPattern)),
				goqu.C(COLUMN_SLUG).Eq(options.Search),
			)
		default:
			// PostgreSQL, MySQL: use ILike
			searchExpr = goqu.Or(
				goqu.C(COLUMN_NAME).ILike("%"+options.Search+"%"),
				goqu.C(COLUMN_DESCRIPTION).ILike("%"+options.Search+"%"),
				goqu.C(COLUMN_SLUG).Eq(options.Search),
			)
		}
		q = q.Where(searchExpr)
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}

		sortOrder := "asc"
		if options.SortOrder != "" {
			sortOrder = options.SortOrder
		}

		orderBy := COLUMN_NAME
		if options.OrderBy != "" {
			orderBy = options.OrderBy
		}

		if sortOrder == sb.ASC {
			q = q.Order(goqu.I(orderBy).Asc())
		} else {
			q = q.Order(goqu.I(orderBy).Desc())
		}
	}

	return q
}

// ============================ POST-TERM RELATIONSHIP METHODS ============================

// PostInsertTermAt creates a relationship between a post and a term at a specific sequence position.
// Also increments the term's count. Duplicate key errors are ignored.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) PostInsertTermAt(ctx context.Context, postID string, termID string, sequence int) error {
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if postID == "" || termID == "" {
		return errors.New("post id and term id are required")
	}

	relation := NewTermRelation()
	relation.SetPostID(postID).
		SetTermID(termID).
		SetSequence(sequence)

	data := relation.GetData()

	sqlStr, sqlParams, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.termRelationTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, sqlParams...)
	if err != nil {
		// Ignore duplicate key errors
		if isDuplicateKeyError(err) {
			return nil
		}
		return err
	}

	// Increment term count
	return store.TermIncrementCount(ctx, termID)
}

// PostRemoveTerm removes the relationship between a post and a term.
// Also decrements the term's count.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) PostRemoveTerm(ctx context.Context, postID string, termID string) error {
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if postID == "" || termID == "" {
		return errors.New("post id and term id are required")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.termRelationTableName).
		Where(
			goqu.C(COLUMN_POST_ID).Eq(postID),
			goqu.C(COLUMN_TERM_ID).Eq(termID),
		).
		Prepared(true).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)
	if err != nil {
		return err
	}

	// Decrement term count
	return store.TermDecrementCount(ctx, termID)
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

	// Get all term relations for this post
	relations, err := store.termRelationList(ctx, TermRelationQueryOptions{PostID: postID})
	if err != nil {
		return []TermInterface{}, err
	}

	if len(relations) == 0 {
		return []TermInterface{}, nil
	}

	// Get term IDs
	termIDs := lo.Map(relations, func(r TermRelationInterface, _ int) string {
		return r.GetTermID()
	})

	// Get terms
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

		terms = lo.Filter(terms, func(t TermInterface, _ int) bool {
			return t.GetTaxonomyID() == taxonomy.GetID()
		})
	}

	return terms, nil
}

// PostListByTermID retrieves all posts associated with a specific term.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) PostListByTermID(ctx context.Context, termID string, options PostQueryOptions) ([]PostInterface, error) {
	if !store.taxonomyEnabled {
		return []PostInterface{}, errors.New("taxonomy is not enabled")
	}
	if termID == "" {
		return []PostInterface{}, errors.New("term id is required")
	}

	// Get all term relations for this term
	relations, err := store.termRelationList(ctx, TermRelationQueryOptions{TermID: termID})
	if err != nil {
		return []PostInterface{}, err
	}

	if len(relations) == 0 {
		return []PostInterface{}, nil
	}

	// Get post IDs
	postIDs := lo.Map(relations, func(r TermRelationInterface, _ int) string {
		return r.GetPostID()
	})

	// Get posts
	options.IDIn = postIDs
	return store.PostList(ctx, options)
}

// PostSetTerms sets the terms for a post within a specific taxonomy.
// Removes any existing terms not in the provided list and adds new ones.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) PostSetTerms(ctx context.Context, postID string, taxonomySlug string, termIDs []string) error {
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if postID == "" {
		return errors.New("post id is required")
	}

	// Get taxonomy
	var taxonomyID string
	if taxonomySlug != "" {
		taxonomy, err := store.TaxonomyFindBySlug(ctx, taxonomySlug)
		if err != nil {
			return err
		}
		if taxonomy == nil {
			return errors.New("taxonomy not found: " + taxonomySlug)
		}
		taxonomyID = taxonomy.GetID()
	}

	// Get current term relations for this post and taxonomy
	currentTerms, err := store.TermListByPostID(ctx, postID, taxonomySlug)
	if err != nil {
		return err
	}

	// Remove terms that are no longer in the list
	for _, currentTerm := range currentTerms {
		found := false
		for _, newTermID := range termIDs {
			if currentTerm.GetID() == newTermID {
				found = true
				break
			}
		}
		if !found {
			if err := store.PostRemoveTerm(ctx, postID, currentTerm.GetID()); err != nil {
				return err
			}
		}
	}

	// Add new terms
	for i, termID := range termIDs {
		if termID == "" {
			continue
		}

		// Verify term exists and belongs to the taxonomy
		if taxonomyID != "" {
			term, err := store.TermFindByID(ctx, termID)
			if err != nil {
				return err
			}
			if term == nil || term.GetTaxonomyID() != taxonomyID {
				continue
			}
		}

		if err := store.PostInsertTermAt(ctx, postID, termID, i); err != nil {
			return err
		}
	}

	return nil
}

// termRelationList retrieves term relations based on query options.
func (store *storeImplementation) termRelationList(ctx context.Context, options TermRelationQueryOptions) ([]TermRelationInterface, error) {
	q := goqu.Dialect(store.dbDriverName).
		From(store.termRelationTableName)

	if options.PostID != "" {
		q = q.Where(goqu.C(COLUMN_POST_ID).Eq(options.PostID))
	}

	if options.TermID != "" {
		q = q.Where(goqu.C(COLUMN_TERM_ID).Eq(options.TermID))
	}

	sqlStr, sqlParams, errSql := q.Select().
		Prepared(true).
		ToSQL()

	if errSql != nil {
		return []TermRelationInterface{}, errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(
		database.NewQueryableContext(ctx, store.db),
		sqlStr,
		sqlParams...,
	)
	if err != nil {
		return []TermRelationInterface{}, err
	}

	list := []TermRelationInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewTermRelationFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

// ============================ UTILITY METHODS ============================

// TermIncrementCount increments the count of posts associated with a term.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermIncrementCount(ctx context.Context, termID string) error {
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if termID == "" {
		return errors.New("term id is required")
	}

	sqlStr := "UPDATE " + store.termTableName + " SET " + COLUMN_COUNT + " = " + COLUMN_COUNT + " + 1 WHERE " + COLUMN_ID + " = ?"

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, termID)
	return err
}

// TermDecrementCount decrements the count of posts associated with a term.
// The count will not go below zero.
// Returns an error if taxonomy features are not enabled.
func (store *storeImplementation) TermDecrementCount(ctx context.Context, termID string) error {
	if !store.taxonomyEnabled {
		return errors.New("taxonomy is not enabled")
	}
	if termID == "" {
		return errors.New("term id is required")
	}

	sqlStr := "UPDATE " + store.termTableName + " SET " + COLUMN_COUNT + " = MAX(0, " + COLUMN_COUNT + " - 1) WHERE " + COLUMN_ID + " = ?"

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, termID)
	return err
}

// isDuplicateKeyError checks if the error is a duplicate key/unique constraint violation.
// It supports SQLite, MySQL, and PostgreSQL error patterns.
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// Check for common duplicate key error patterns
	return contains(errStr, "duplicate") ||
		contains(errStr, "UNIQUE constraint failed") ||
		contains(errStr, "1062") || // MySQL duplicate entry
		contains(errStr, "23505") // PostgreSQL unique violation
}

// contains checks if string s contains the substring substr.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

// containsAt checks if string s contains the substring substr starting at position start.
func containsAt(s, substr string, start int) bool {
	if start+len(substr) > len(s) {
		return false
	}
	for i := 0; i < len(substr); i++ {
		if s[start+i] != substr[i] {
			return containsAt(s, substr, start+1)
		}
	}
	return true
}
