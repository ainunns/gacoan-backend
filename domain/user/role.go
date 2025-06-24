package user

import "fmt"

const (
	RoleSuperAdmin = "superadmin"
	RoleCustomer   = "customer"
	RoleKitchen    = "kitchen"
	RoleWaiter     = "waiter"
)

var (
	Roles = []Role{
		{RoleSuperAdmin},
		{RoleCustomer},
		{RoleKitchen},
		{RoleWaiter},
	}
)

type Role struct {
	Name string
}

func NewRole(name string) (Role, error) {
	if !isValidRole(name) {
		return Role{}, fmt.Errorf("invalid role Name")
	}
	return Role{
		Name: name,
	}, nil
}

func NewRoleFromSchema(name string) Role {
	return Role{
		Name: name,
	}
}

func isValidRole(name string) bool {
	for _, role := range Roles {
		if role.Name == name {
			return true
		}
	}
	return false
}
