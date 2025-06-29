package transaction

import (
	"fmt"
	"strings"
)

type QueueCode struct {
	Code  string
	Valid bool
}

func NewQueueCode(code string) (QueueCode, error) {
	return QueueCode{
		Code:  code,
		Valid: strings.HasPrefix(code, "Q") && len(code) == 5 && code[1:] != "0000",
	}, nil
}

func NewQueueCodeFromSchema(code string, valid bool) QueueCode {
	return QueueCode{
		Code:  code,
		Valid: valid,
	}
}

func (q *QueueCode) QueueNumber() (int, error) {
	number := strings.TrimPrefix(q.Code, "Q")
	if number == "0000" || number == "" {
		return 0, nil
	}
	var queueNumber int
	_, err := fmt.Sscanf(number, "%4d", &queueNumber)
	if err != nil {
		return 0, fmt.Errorf("invalid queue code: %s", q.Code)
	}
	return queueNumber, nil
}
