package service

import (
	"context"
	"errors"
	"fp-kpl/application/response"
	"fp-kpl/domain/table"

	"gorm.io/gorm"
)

type (
	TableService interface {
		GetAllTables(ctx context.Context) ([]response.Table, error)
		GetTableByID(ctx context.Context, id string) (response.Table, error)
	}
)

type tableService struct {
	tableRepository table.Repository
}

func NewTableService(tableRepository table.Repository) TableService {
	return &tableService{tableRepository: tableRepository}
}

func (s *tableService) GetAllTables(ctx context.Context) ([]response.Table, error) {
	retrievedTables, err := s.tableRepository.GetAllTables(ctx, nil)
	if err != nil {
		return nil, table.ErrorGetAllTables
	}

	responseTables := make([]response.Table, 0, len(retrievedTables))
	for _, table := range retrievedTables {
		responseTables = append(responseTables, response.Table{
			ID:          table.ID.String(),
			TableNumber: table.TableNumber,
		})
	}
	return responseTables, nil
}

func (s *tableService) GetTableByID(ctx context.Context, id string) (response.Table, error) {
	retrievedTable, err := s.tableRepository.GetTableByID(ctx, nil, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Table{}, table.ErrorTableNotFound
		}
		return response.Table{}, table.ErrorGetTableByID
	}

	return response.Table{
		ID:          retrievedTable.ID.String(),
		TableNumber: retrievedTable.TableNumber,
	}, nil
}
