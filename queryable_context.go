package blogstore

import (
	"context"

	"github.com/dracory/database"
)

// toQueryableContext converts a context.Context to a database.QueryableContext.
// If the context is already a QueryableContext, it returns it directly.
// Otherwise, it wraps the context with the store's database connection.
func (store *storeImplementation) toQueryableContext(ctx context.Context) database.QueryableContext {
	if database.IsQueryableContext(ctx) {
		return ctx.(database.QueryableContext)
	}
	return database.NewQueryableContext(ctx, store.db)
}
