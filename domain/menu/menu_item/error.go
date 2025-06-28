package menu

import "errors"

var (
	ErrorGetAllMenus            = errors.New("failed to get all menus")
	ErrorGetMenuByID            = errors.New("failed to get menu by id")
	ErrorGetMenusByCategoryID   = errors.New("failed to get menus by category id")
	ErrorCategoryNotFound       = errors.New("category not found")
	ErrorMenuNotFound           = errors.New("menu not found")
	ErrorUpdateMenuAvailability = errors.New("failed to update menu availability")
)
