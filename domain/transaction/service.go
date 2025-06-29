package transaction

import (
	"context"

	"fmt"
)

type (
	Service interface {
		GenerateQueueCode(ctx context.Context, transactionID string) (string, error)
	}

	service struct {
		transactionRepository Repository
	}
)

func NewService(transactionRepository Repository) Service {
	return &service{
		transactionRepository: transactionRepository,
	}
}

func (s *service) GenerateQueueCode(ctx context.Context, transactionID string) (string, error) {
	latestCode, err := s.transactionRepository.GetLatestQueueCode(ctx, nil, transactionID)
	if err != nil {
		return "", fmt.Errorf("failed to get latest queue code: %w", err)
	}

	return latestCode, nil
}
