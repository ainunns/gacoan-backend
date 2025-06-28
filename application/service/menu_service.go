package service

import (
	"context"
	"fp-kpl/application/response"
	"fp-kpl/domain/menu/category"
	menu "fp-kpl/domain/menu/menu_item"
)

type (
	MenuService interface {
		GetAllMenus(ctx context.Context) ([]response.Menu, error)
		GetMenuByID(ctx context.Context, id string) (response.Menu, error)
		GetMenusByCategoryID(ctx context.Context, categoryID string) ([]response.Menu, error)
	}

	menuService struct {
		menuRepository     menu.Repository
		categoryRepository category.Repository
	}
)

func NewMenuService(menuRepository menu.Repository, categoryRepository category.Repository) MenuService {
	return &menuService{menuRepository: menuRepository, categoryRepository: categoryRepository}
}

func (s *menuService) GetAllMenus(ctx context.Context) ([]response.Menu, error) {
	retrievedMenus, err := s.menuRepository.GetAllMenus(ctx, nil)
	if err != nil {
		return nil, menu.ErrorGetAllMenus
	}

	responseMenus := make([]response.Menu, 0, len(retrievedMenus))
	for _, menu := range retrievedMenus {
		category, _ := s.categoryRepository.GetCategoryByID(ctx, nil, menu.CategoryID.String())

		responseMenus = append(responseMenus, response.Menu{
			ID:          menu.ID.String(),
			Name:        menu.Name,
			Description: menu.Description,
			Price:       menu.Price.Price,
			Category: response.Category{
				ID:   category.ID.String(),
				Name: category.Name,
			},
		})
	}

	return responseMenus, nil
}

func (s *menuService) GetMenuByID(ctx context.Context, id string) (response.Menu, error) {
	retrievedMenu, err := s.menuRepository.GetMenuByID(ctx, nil, id)
	if err != nil {
		return response.Menu{}, menu.ErrorGetMenuByID
	}

	category, _ := s.categoryRepository.GetCategoryByID(ctx, nil, retrievedMenu.CategoryID.String())

	responseMenu := response.Menu{
		ID:          retrievedMenu.ID.String(),
		Name:        retrievedMenu.Name,
		Description: retrievedMenu.Description,
		Price:       retrievedMenu.Price.Price,
		Category: response.Category{
			ID:   category.ID.String(),
			Name: category.Name,
		},
	}

	return responseMenu, nil
}

func (s *menuService) GetMenusByCategoryID(ctx context.Context, categoryID string) ([]response.Menu, error) {
	retrievedMenus, err := s.menuRepository.GetMenusByCategoryID(ctx, nil, categoryID)
	if err != nil {
		return nil, menu.ErrorGetAllMenus
	}

	responseMenus := make([]response.Menu, 0, len(retrievedMenus))
	for _, menu := range retrievedMenus {
		category, _ := s.categoryRepository.GetCategoryByID(ctx, nil, menu.CategoryID.String())

		responseMenus = append(responseMenus, response.Menu{
			ID:          menu.ID.String(),
			Name:        menu.Name,
			Description: menu.Description,
			Price:       menu.Price.Price,
			Category: response.Category{
				ID:   category.ID.String(),
				Name: category.Name,
			},
		})
	}

	return responseMenus, nil
}
