package database

import (
	"errors"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/google/uuid"
)

func (d *Database) Logout(userId uuid.UUID) error {
	sql := "DELETE FROM logged_in_user WHERE user_id = $1"

	_, err := d.Exec(sql, userId)
	return err
}

func (d *Database) Login(emailPar, passwordPar string) (string, error) {
	var id, email, name, password, company_id, role_id string
	sql := "SELECT id, email, name, password, company_id, role_id from \"user\" where email = $1"
	if err := d.QueryRow(sql, emailPar).Scan(&id, &email, &name, &password, &company_id, &role_id); err != nil {
		err := errors.New("Credenciales no válidas")
		return "", err
	}

	if _, err := utils.ComparePassword(password, passwordPar); err != nil {
		err := errors.New("Credenciales no válidas")
		return "", err
	}
	uId, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	cId, err := uuid.Parse(company_id)
	if err != nil {
		return "", err
	}

	u := types.User{
		Id:        uId,
		Email:     email,
		Name:      name,
		CompanyId: cId,
		RoleId:    role_id,
	}
	token, err := utils.GenerateToken(u)
	sql = "INSERT INTO logged_in_user (user_id, token) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET token = $2"
	if _, err := d.Exec(sql, u.Id, []byte(token)); err != nil {
		return "", err
	}

	return token, nil
}
