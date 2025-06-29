package port

type (
	HashPort interface {
		GenerateHash(input string) (string, error)
	}
)
