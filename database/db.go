package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	if _, err := loadRoles(db); err != nil {
		log.Println(err)
	}

	return db, nil
}
