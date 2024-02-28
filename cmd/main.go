package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("error to load .env file: [%s]\n", err.Error())
	}

	args := os.Args[1:]
	if len(args) <= 0 {
		log.Print("you need make at least one order\n")
	}

	db, err := runPostgres()
	if err != nil {
		log.Println(err)
	}

	if err := initQueries(db); err != nil {
		log.Println(err)
	}

	repository := NewRepository(db)
	service := NewService(repository)

	orders, err := service.execute(args)
	if err != nil {
		log.Println(err)
	}

	for key, val := range orders {
		fmt.Printf("===Стелаж %s\n", key)
		for _, k := range val {
			if len(k.OtherShelfs) != 0 {
				fmt.Printf("%s (id=%d) \nзаказ %d, %d шт\nдоп стелаж: %v\n", k.ItemName, k.ItemId, k.OrderId, k.Quantity, string(k.OtherShelfs))
				fmt.Println()
				continue
			}
			fmt.Printf("%s (id=%d) \nзаказ %d, %d шт\n", k.ItemName, k.ItemId, k.OrderId, k.Quantity)
			fmt.Println()
		}
	}
}

func initQueries(db *sql.DB) error {
	file, err := os.ReadFile("sql/query.sql")
	if err != nil {
		return fmt.Errorf("error to read the given file: [%s]", err.Error())
	}

	rows := strings.Split(string(file), ";")
	for _, row := range rows {
		if _, err = db.Exec(row); err != nil {
			fmt.Printf("error to exec the row: [%s]\n", err.Error())
		}
	}

	return nil
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
