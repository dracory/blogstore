package blogstore

import (
	"context"

	"github.com/dracory/database"
)

func (store *store) toQueryableContext(ctx context.Context) database.QueryableContext {
	if database.IsQueryableContext(ctx) {
		return ctx.(database.QueryableContext)
	}
	return database.NewQueryableContext(ctx, store.db)
}
