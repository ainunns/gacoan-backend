package transaction

import (
	"fmt"
	"strings"
)

type QueueCode struct {
	Code string
}

func NewQueueCode(code string) (QueueCode, error) {
	return QueueCode{
		Code: code,
	}, nil
}

func (q *QueueCode) QueueNumber() (int, error) {
	number := strings.TrimPrefix(q.Code, "Q")
	if number == "0000" {
		return 0, nil
	}
	var queueNumber int
	_, err := fmt.Sscanf(number, "%4d", &queueNumber)
	if err != nil {
		return 0, fmt.Errorf("invalid queue code: %s", q.Code)
	}
	return queueNumber, nil
}
