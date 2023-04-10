package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Connect() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := sql.Open("mysql", os.Getenv("DB_ROLE")+":"+os.Getenv("DB_PASSWORD")+"@tcp(localhost:3306)/"+os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
