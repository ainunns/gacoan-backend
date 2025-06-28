package menu

import "context"

type (
	Repository interface {
		GetAllMenus(ctx context.Context, tx interface{}) ([]Menu, error)
		GetMenuByID(ctx context.Context, tx interface{}, id string) (Menu, error)
		GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]Menu, error)
	}
)
