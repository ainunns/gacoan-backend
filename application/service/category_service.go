package service

import (
	"context"
	"fp-kpl/application/response"
	"fp-kpl/domain/menu/category"
)

type (
	CategoryService interface {
		GetAllCategories(ctx context.Context) ([]response.Category, error)
		GetCategoryByID(ctx context.Context, id string) (response.Category, error)
	}

	categoryService struct {
		categoryRepository category.Repository
	}
)

func NewCategoryService(categoryRepository category.Repository) CategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]response.Category, error) {
	retrievedCategories, err := s.categoryRepository.GetAllCategories(ctx, nil)
	if err != nil {
		return nil, category.ErrorGetAllCategories
	}

	responseCategories := make([]response.Category, 0, len(retrievedCategories))
	for _, category := range retrievedCategories {
		responseCategories = append(responseCategories, response.Category{
			ID:   category.ID.String(),
			Name: category.Name,
		})
	}

	return responseCategories, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (response.Category, error) {
	retrievedCategory, err := s.categoryRepository.GetCategoryByID(ctx, nil, id)
	if err != nil {
		return response.Category{}, category.ErrorGetCategoryByID
	}

	return response.Category{
		ID:   retrievedCategory.ID.String(),
		Name: retrievedCategory.Name,
	}, nil
}
