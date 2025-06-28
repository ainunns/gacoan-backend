package category

import "errors"

var (
	ErrorGetAllCategories = errors.New("failed to get all categories")
	ErrorGetCategoryByID  = errors.New("failed to get category by id")
	ErrorCategoryNotFound = errors.New("category not found")
)
