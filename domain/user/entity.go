package user

import (
	"fp-kpl/domain/identity"
	"fp-kpl/domain/shared"
)

type User struct {
	ID          identity.ID
	Email       string
	Password    Password
	Name        string
	PhoneNumber string
	Role        Role
	shared.Timestamp
}
