package category

import "context"

type (
	Repository interface {
		GetAllCategory(ctx context.Context, tx interface{}) ([]Category, error)
		GetCategoryByID(ctx context.Context, tx interface{}, id string) (Category, error)
	}
)
