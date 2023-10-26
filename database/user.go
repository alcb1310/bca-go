package database

import (
	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/google/uuid"
)

type UserInfo struct {
	Email     *string
	Name      *string
	Password  *string
	Role      *string
	CompanyId uuid.UUID
}

func (d *Database) AddUser(u UserInfo) error {
	var pass []byte
	var err error
	if pass, err = utils.EncryptPasssword(*u.Password); err != nil {
		return err
	}
	password := string(pass)
	u.Password = &password

	sql := "INSERT INTO \"user\" (email, name, password, company_id, role_id) VALUES($1, $2, $3, $4, $5)"
	if _, err := d.Exec(sql, &u.Email, &u.Name, &u.Password, u.CompanyId, &u.Role); err != nil {
		return err
	}

	return nil
}

func (d *Database) GetAllUsers(company_id uuid.UUID) ([]types.User, error) {
	sql := "SELECT user_id, user_email, user_name, role_name FROM user_without_password WHERE company_id=$1"
	rows, err := d.Query(sql, company_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []types.User
	for rows.Next() {
		var id, email, name, role string
		if err := rows.Scan(&id, &email, &name, &role); err != nil {
			return nil, err
		}

		strUUID, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		users = append(users, types.User{
			Id:     strUUID,
			Email:  email,
			Name:   name,
			RoleId: role,
		})
	}

	return users, nil
}
