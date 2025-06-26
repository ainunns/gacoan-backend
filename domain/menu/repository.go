package menu

import "context"

type (
	Repository interface {
		GetAllMenu(ctx context.Context, tx interface{}) ([]Menu, error)
		GetMenuByID(ctx context.Context, tx interface{}, id string) (Menu, error)
		GetMenuByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]Menu, error)
	}
)
