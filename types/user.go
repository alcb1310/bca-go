package types

import "github.com/google/uuid"

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CompanyId uuid.UUID `json:"company_id"`
	RoleId    string    `json:"role_id"`
}
