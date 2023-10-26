package database

import (
	"strings"

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
	_, err = d.Exec(sql, &u.Email, &u.Name, &u.Password, u.CompanyId, &u.Role)
	return err
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

func (d *Database) DeleteUser(user_id, company_id uuid.UUID) error {
	sql := "DELETE FROM \"user\" WHERE id = $1 and company_id = $2"
	_, err := d.Exec(sql, user_id, company_id)
	return err
}

func (d *Database) GetOneUser(user_id, company_id uuid.UUID) (types.User, error) {
	u := types.User{}

	sql := "SELECT user_id, user_email, user_name, role_id FROM user_without_password where company_id = $2 and user_id = $1"
	rows, err := d.Query(sql, user_id, company_id)
	if err != nil {
		return u, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, email, name, roleId string
		if err := rows.Scan(&id, &email, &name, &roleId); err != nil {
			return u, err
		}
		foundUserId, err := uuid.Parse(id)
		if err != nil {
			return u, err
		}

		u.Id = foundUserId
		u.Name = name
		u.Email = email
		u.CompanyId = company_id
		u.RoleId = strings.Trim(roleId, " ")
	}

	return u, nil
}

func (d *Database) UpdateUser(u *UserInfo, user_id, company_id uuid.UUID) error {
	sql := "UPDATE \"user\" SET email=$3, name=$4, role_id=$5 WHERE id=$1 AND company_id = $2"
	_, err := d.Exec(sql, user_id, company_id, u.Email, u.Name, u.Role)
	return err
}

func (d *Database) ChangePassword(oldPassword, newPassword string, user_id, company_id uuid.UUID) error {
	var prev *string

	sql := "SELECT password FROM \"user\" WHERE id=$1 and company_id = $2"
	rows, err := d.Query(sql, user_id, company_id)
	if err != nil {
		return err
	}
	defer rows.Close()

	prev = nil
	for rows.Next() {
		var passwd string
		if err := rows.Scan(&passwd); err != nil {
			return err
		}

		prev = &passwd
	}

	if _, err := utils.ComparePassword(*prev, oldPassword); err != nil {
		return err
	}

	savePasswd, err := utils.EncryptPasssword(newPassword)
	if err != nil {
		return err
	}

	sql = "UPDATE \"user\" SET password = $3 WHERE id=$1 and company_id = $2"
	_, err = d.Exec(sql, savePasswd, user_id, company_id)

	return err
}
