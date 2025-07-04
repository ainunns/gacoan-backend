package transaction

import (
	"context"
	"time"

	"fmt"
)

type (
	Service interface {
		GenerateQueueCode(ctx context.Context, transactionID string) (string, error)
		CalculateMaxCookingTime(orders []OrderQuery) time.Duration
		GetOrderDelayStatus(maxCookingTime time.Duration, cookedAt *time.Time, servedAt *time.Time) bool
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

func (s *service) CalculateMaxCookingTime(orders []OrderQuery) time.Duration {
	maxCookingTime := time.Duration(0)

	for _, orderQuery := range orders {
		if orderQuery.Menu.CookingTime > maxCookingTime {
			maxCookingTime = orderQuery.Menu.CookingTime
		}
	}

	return maxCookingTime
}

func (s *service) GetOrderDelayStatus(maxCookingTime time.Duration, cookedAt *time.Time, servedAt *time.Time) bool {
	now := time.Now()
	isDelayed := false
	if cookedAt != nil {
		expectedFinishTime := cookedAt.Add(maxCookingTime)
		if servedAt != nil {
			isDelayed = servedAt.After(expectedFinishTime)
		} else {
			isDelayed = now.After(expectedFinishTime)
		}
	}

	return isDelayed
}
