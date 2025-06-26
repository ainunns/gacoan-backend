package table

import "context"

type (
	Repository interface {
		GetAllTables(ctx context.Context, tx interface{}) ([]Table, error)
		GetTableByID(ctx context.Context, tx interface{}, id string) (Table, error)
	}
)
