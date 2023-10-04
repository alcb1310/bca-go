package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	database := os.Getenv("PGDATABASE")
	password := os.Getenv("PGPASSWORD")
	user := os.Getenv("PGUSER")

	dbString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		log.Println(err)
	}
	return db, nil
}

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
