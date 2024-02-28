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
				fmt.Printf("%s (id=%d) \nзаказ %d, %d шт\nдоп стелаж: %v\n", k.ItemName, k.ItemId, k.OrderId, k.Quantity, k.OtherShelfs)
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

// func (p *Postgres) getOrder(args []string) (*Order, error) {
// 	var orderDesc OrderDescription
// 	order := Order{
// 		Id:               1,
// 		OrderDescription: []OrderDescription{},
// 	}

// 	for _, arg := range args {
// 		numArg, err := strconv.Atoi(arg)
// 		if err != nil {
// 			fmt.Printf("error to convert arg to int type: [%s]\n", err.Error())
// 			continue
// 		}

// 		rows, err := p.db.Query("SELECT * FROM order_description WHERE id = $1", numArg)
// 		if err != nil {
// 			return nil, fmt.Errorf("error to select order_description: [%s]", err.Error())
// 		}
// 		defer rows.Close()

// 		for rows.Next() {
// 			err := rows.Scan(&orderDesc.Id, &orderDesc.ItemId, &orderDesc.Quantity)
// 			if err != nil {
// 				fmt.Printf("error to scan rows: [%s]\n", err.Error())
// 			}
// 			order.OrderDescription = append(order.OrderDescription, orderDesc)
// 		}
// 		if err := rows.Err(); err != nil {
// 			log.Fatal(err)
// 		}

// 		finalOrders, err := p.getOrderDescription(&order)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		for _, finalOrder := range finalOrders {
// 			fmt.Printf("%+v\n", finalOrder)
// 		}
// 	}

// 	return &order, nil
// }

// func (p *Postgres) getOrderDescription(order *Order) ([]FinalOrder, error) {
// 	var finalOrders []FinalOrder
// 	for _, orderDesc := range order.OrderDescription {
// 		item, err := p.selectItemById(orderDesc.ItemId)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		// Get products from shelfs
// 		quantity, err := p.getItemsFromShelf(item)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		// Create final order
// 		finalOrder := FinalOrder{
// 			Shelf:    string(item.MainShelf),
// 			Products: []Product{},
// 		}
// 		product := Product{
// 			ItemName:    item.ItemName,
// 			ItemId:      item.Id,
// 			OrderNumber: orderDesc.Id,
// 			ItemCount:   quantity,
// 		}

// 		// Start to search in other shelfs if we don't get enough quantity of products
// 		if quantity != orderDesc.Quantity {

// 			if item.OtherShelfs != nil {
// 				// Take items from item other shelfs
// 				quantity, err := p.getItemsFromOtherShelf(item, quantity, orderDesc.Quantity, 100)
// 				if err != nil {
// 					log.Println(err)
// 				}

// 				if quantity != orderDesc.Quantity {
// 					// Make recursive search of items from shelfs
// 					fmt.Println("sorry don't have enough items in stock")
// 					return finalOrders, nil
// 				}

// 			}

// 		}

// 		finalOrder.Products = append(finalOrder.Products, product)
// 		finalOrders = append(finalOrders, finalOrder)
// 	}

// 	return finalOrders, nil
// }

// // Get items from a shelf
// func (p *Postgres) getItemsFromShelf(item *Item) (int64, error) {
// 	var shelf Shelf
// 	var arrIds pq.Int64Array
// 	err := p.db.QueryRow("SELECT * FROM shelfs WHERE shelf_type = $1 LIMIT 1", item.MainShelf).
// 		Scan(&shelf.Id, &shelf.ShelfType, &arrIds)
// 	if err != nil {
// 		return 0, fmt.Errorf("error getItemsFromShelf to scan row: [%s]", err.Error())
// 	}
// 	shelf.Items = arrIds

// 	var quantity int64 = 0
// 	for i, id := range shelf.Items {
// 		if id != item.Id {
// 			continue
// 		}
// 		quantity += 1
// 		shelf.Items[i] = 0
// 		if err := p.updateShelfState(shelf.Id, shelf.Items); err != nil {
// 			log.Println(err)
// 		}
// 	}

// 	return quantity, nil
// }

// // Get items from item other shelfs, if we don't get enough item quantity from shelf
// // and if item other shelf != nil
// func (p *Postgres) getItemsFromOtherShelf(item *Item, count, orderCount, limit int64) (int64, error) {
// 	var shelf Shelf
// 	var arrIds pq.Int64Array
// 	var quantity int64 = count
// 	for _, ty := range item.OtherShelfs {
// 		rows, err := p.db.Query("SELECT * FROM shelfs WHERE shelf_type = $1 LIMIT $2", ty, limit)
// 		if err != nil {
// 			return 0, fmt.Errorf("error getItemsFromOtherShelf to scan row: [%s]", err.Error())
// 		}
// 		defer rows.Close()
// 		shelf.Items = arrIds // Use arrIds instead this

// 		for rows.Next() {
// 			err := rows.Scan(&shelf.Id, &shelf.ShelfType, &arrIds)
// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			for i, id := range shelf.Items {
// 				if id != item.Id {
// 					continue
// 				}
// 				quantity += 1
// 				shelf.Items[i] = 0
// 				if err := p.updateShelfState(shelf.Id, shelf.Items); err != nil {
// 					log.Println(err)
// 				}

// 				if quantity == orderCount {
// 					return quantity, nil
// 				}
// 			}
// 		}
// 		if err := rows.Err(); err != nil {
// 			log.Fatal(err)
// 		}

// 	}

// 	return quantity, nil
// }
