package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func runPostgres() (*sql.DB, error) {
	godotenv.Load()
	var (
		dbuser     = os.Getenv("DB_USER")
		dbname     = os.Getenv("DB_NAME")
		dbpassword = os.Getenv("DB_PASSWORD")
		dbhost     = os.Getenv("DB_HOST")
		dbport     = os.Getenv("DB_PORT")
		uri        = fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s",
			dbuser, dbname, dbpassword, dbhost, dbport)
	)

	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("error to open psql connection: [%s]", err.Error())
	}

	return db, nil
}
func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("error to load .env file: [%s]\n", err.Error())
	}

	db, err := runPostgres()
	if err != nil {
		log.Println(err)
	}

	psql := NewPostgres(db)
	if err := psql.dropAllTables(); err != nil {
		log.Println(err)
	}
}

func (p *Postgres) dropAllTables() error {
	_, err := p.db.Exec(`DROP TABLE items, shelfs, order_description;`)
	if err != nil {
		return fmt.Errorf("error to rop tables: [%s]", err.Error())
	}

	return nil
}
