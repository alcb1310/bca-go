package types

import "github.com/google/uuid"

type CreateCompany struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Employees    uint8  `json:"employees"`
	UserEmail    string `json:"user_email"`
	UserPassword string `json:"user_password"`
	UserName     string `json:"user_name"`
}

type Company struct {
	ID        uuid.UUID `json:"id"`
	RUC       string    `json:"ruc"`
	Name      string    `json:"name"`
	Employees uint8     `json:"employees"`
	IsActive  bool      `json:"is_active"`
}
