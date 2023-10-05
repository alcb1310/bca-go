package types

type CreateCompany struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Employees    uint8  `json:"employees"`
	UserEmail    string `json:"user_email"`
	UserPassword string `json:"user_password"`
	UserName     string `json:"user_name"`
}
