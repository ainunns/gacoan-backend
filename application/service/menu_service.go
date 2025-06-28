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
		UpdateMenuAvailability(ctx context.Context, id string, isAvailable bool) (response.Menu, error)
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
		return nil, menu.ErrorCategoryNotFound
	}

	responseMenus := make([]response.Menu, 0, len(retrievedMenus))
	for _, menu := range retrievedMenus {
		categoryDetail, err := s.categoryRepository.GetCategoryByID(ctx, nil, menu.CategoryID.String())

		if err != nil {
			return nil, category.ErrorGetCategoryByID
		}

		responseMenus = append(responseMenus, response.Menu{
			ID:          menu.ID.String(),
			Name:        menu.Name,
			Description: menu.Description,
			ImageUrl:    menu.ImageURL.Path,
			IsAvailable: menu.IsAvailable,
			Price:       menu.Price.Price,
			Category: response.Category{
				ID:   categoryDetail.ID.String(),
				Name: categoryDetail.Name,
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

	categoryDetail, err := s.categoryRepository.GetCategoryByID(ctx, nil, retrievedMenu.CategoryID.String())
	if err != nil {
		return response.Menu{}, category.ErrorGetCategoryByID
	}

	responseMenu := response.Menu{
		ID:          retrievedMenu.ID.String(),
		Name:        retrievedMenu.Name,
		Description: retrievedMenu.Description,
		ImageUrl:    retrievedMenu.ImageURL.Path,
		IsAvailable: retrievedMenu.IsAvailable,
		Price:       retrievedMenu.Price.Price,
		Category: response.Category{
			ID:   categoryDetail.ID.String(),
			Name: categoryDetail.Name,
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
		categoryDetail, err := s.categoryRepository.GetCategoryByID(ctx, nil, menu.CategoryID.String())

		if err != nil {
			return nil, category.ErrorGetCategoryByID
		}

		responseMenus = append(responseMenus, response.Menu{
			ID:          menu.ID.String(),
			Name:        menu.Name,
			Description: menu.Description,
			ImageUrl:    menu.ImageURL.Path,
			IsAvailable: menu.IsAvailable,
			Price:       menu.Price.Price,
			Category: response.Category{
				ID:   categoryDetail.ID.String(),
				Name: categoryDetail.Name,
			},
		})
	}

	return responseMenus, nil
}

func (s *menuService) UpdateMenuAvailability(ctx context.Context, id string, isAvailable bool) (response.Menu, error) {
	updatedMenu, err := s.menuRepository.UpdateMenuAvailability(ctx, nil, id, isAvailable)
	if err != nil {
		return response.Menu{}, menu.ErrorUpdateMenuAvailability
	}

	categoryDetail, err := s.categoryRepository.GetCategoryByID(ctx, nil, updatedMenu.CategoryID.String())
	if err != nil {
		return response.Menu{}, category.ErrorGetCategoryByID
	}

	responseMenu := response.Menu{
		ID:          updatedMenu.ID.String(),
		Name:        updatedMenu.Name,
		Description: updatedMenu.Description,
		ImageUrl:    updatedMenu.ImageURL.Path,
		IsAvailable: updatedMenu.IsAvailable,
		Price:       updatedMenu.Price.Price,
		Category: response.Category{
			ID:   categoryDetail.ID.String(),
			Name: categoryDetail.Name,
		},
	}

	return responseMenu, nil
}
