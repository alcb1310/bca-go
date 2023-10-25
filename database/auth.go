package database

import "github.com/google/uuid"

func (d *Database) Logout(userId uuid.UUID) error {
	sql := "DELETE FROM logged_in_user WHERE user_id = $1"

	_, err := d.Exec(sql, userId)
	return err
}
