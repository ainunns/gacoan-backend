package table

import "context"

type (
	Repository interface {
		GetAllTable(ctx context.Context, tx interface{}) ([]Table, error)
		GetTableByID(ctx context.Context, tx interface{}, id string) (Table, error)
	}
)
