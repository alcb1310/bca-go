package database

import (
	"database/sql"
	"log"
	"os"
	"strings"
)

func createTables(d *sql.DB) error {
	file, err := os.ReadFile("./database/tables.sql")
	if err != nil {
		log.Println(":Error: Couldn't load sql file:", err.Error())
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		if _, err := tx.Exec(request); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func loadRoles(d *sql.DB) (int, error) {
	var c int

	sql := "select count(*) roles from role"
	err := d.QueryRow(sql).Scan(&c)
	if err != nil {
		return 0, err
	}

	if c == 0 {
		tx, err := d.Begin()
		sql = "insert into role (id, name) values ('a', 'admin'), ('u', 'user')"
		_, err = tx.Exec(sql)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		tx.Commit()
		d.QueryRow(sql).Scan(&c)
	}
	return c, nil
}
