package data

import (
	"fp-kpl/domain/user"
	"fp-kpl/infrastructure/database/schema"
	"log"
)

var Users = []schema.User{
	{
		Name:        "Super Admin",
		Email:       "superadmin@gmail.com",
		Password:    getPassword("Superadmin123!"),
		PhoneNumber: "081234567890",
		Role:        "superadmin",
	},
	{
		Name:        "Kitchen",
		Email:       "kitchen@gmail.com",
		Password:    getPassword("Kitchen123!"),
		PhoneNumber: "081234567890",
		Role:        "kitchen",
	},
	{
		Name:        "Chef",
		Email:       "chef@gmail.com",
		Password:    getPassword("Chef123!"),
		PhoneNumber: "081234567890",
		Role:        "kitchen",
	},
	{
		Name:        "Waiter",
		Email:       "waiter@gmail.com",
		Password:    getPassword("Waiter123!"),
		PhoneNumber: "081234567890",
		Role:        "waiter",
	},
	{
		Name:        "Waiter 2",
		Email:       "waiter2@gmail.com",
		Password:    getPassword("Waiter2123!"),
		PhoneNumber: "081234567890",
		Role:        "waiter",
	},
}

func getPassword(password string) string {
	hashedPassword, err := user.NewPassword(password)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}
	return hashedPassword.Password
}
