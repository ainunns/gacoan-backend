package transaction

import (
	"context"
	"fmt"
)

type (
	Service interface {
		GenerateQueueCode(ctx context.Context) (string, error)
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

func (s *service) GenerateQueueCode(ctx context.Context) (string, error) {
	latestCode, err := s.transactionRepository.GetLatestQueueCode(ctx, nil)
	if err != nil {
		return "", err
	}

	numericLatestCode, err := latestCode.QueueNumber()
	if err != nil {
		return "", err
	}

	newCodeString := fmt.Sprintf("Q%04d", numericLatestCode+1)
	return newCodeString, nil
}
