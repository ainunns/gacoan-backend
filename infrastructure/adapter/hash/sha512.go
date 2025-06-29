package hash

import (
	"crypto/sha512"
	"fmt"
	"fp-kpl/domain/port"
)

type sha512Adapter struct{}

func NewSha512Adapter() port.HashPort {
	return &sha512Adapter{}
}

func (s sha512Adapter) GenerateHash(input string) (string, error) {
	hasher := sha512.New()
	hasher.Write([]byte(input))
	generatedSignature := fmt.Sprintf("%x", hasher.Sum(nil))

	return generatedSignature, nil
}
