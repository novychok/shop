package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Shelf struct {
	Id        int64
	ShelfType []uint8
	Items     []int64
}

type Item struct {
	Id          int64
	ItemName    string
	MainShelf   []uint8
	OtherShelfs []uint8
}

type Order struct {
	Id               int64
	OrderDescription []OrderDescription
}

type OrderDescription struct {
	Id       int64
	ItemId   int64
	Quantity int64
}

type FinalOrder struct {
	Shelf    []uint8
	Products []Product
}

type Product struct {
	ItemName    string
	ItemId      int64
	OrderNumber int64
	ItemCount   int64
}

func (p *Postgres) getOrder(args []string) (*Order, error) {
	var orderDesc OrderDescription
	order := Order{
		Id:               1,
		OrderDescription: []OrderDescription{},
	}

	for _, arg := range args {
		numArg, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("error to convert arg to int type: [%s]\n", err.Error())
			continue
		}

		rows, err := p.db.Query("SELECT * FROM order_description WHERE id = $1", numArg)
		if err != nil {
			return nil, fmt.Errorf("error to select order_description: [%s]", err.Error())
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&orderDesc.Id, &orderDesc.ItemId, &orderDesc.Quantity)
			if err != nil {
				fmt.Printf("error to scan rows: [%s]\n", err.Error())
			}
			order.OrderDescription = append(order.OrderDescription, orderDesc)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return &order, nil
}

func (p *Postgres) getOrderDescription(order *Order) error {
	for _, orderDesc := range order.OrderDescription {
		item, err := p.selectItemById(orderDesc.ItemId)
		if err != nil {
			log.Println(err)
		}

		_, err = p.getItemsFromShelf(item, orderDesc.Quantity)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (p *Postgres) getItemsFromShelf(item *Item, quantity int64) ([]FinalOrder, error) {
	var shelf Shelf
	rows, err := p.db.Query("SELECT * FROM shelfs WHERE shelf_type = $1 LIMIT 1", item.MainShelf)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var finalOrders []FinalOrder

	for rows.Next() {
		var arrIds pq.Int64Array
		err := rows.Scan(&shelf.Id, &shelf.ShelfType, &arrIds)
		if err != nil {
			fmt.Printf("error to scan rows: [%s]\n", err.Error())
		}
		shelf.Items = arrIds

		finalOrder := FinalOrder{}

		for i, id := range shelf.Items {

			if finalOrder.ItemCount != quantity && item.Id == id {
				finalOrder.ItemCount += 1
				shelf.Items[i] = 0
			}

		}

		// Update the shelf array state / because we get items from it
		if err := p.updateShelfState(shelf.Id, shelf.Items); err != nil {
			log.Println(err)
		}

		if finalOrder.ItemCount != quantity {
			if err := p.getWithMoreLimit(item, 10, finalOrder, quantity); err != nil {
				log.Println(err)
			}
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return finalOrders, nil
}

func (p *Postgres) getWithMoreLimit(item *Item, limit int64, finalOrder FinalOrder, quantity int64) error {
	var shelf Shelf
	var shelfs []Shelf
	rows, err := p.db.Query("SELECT * FROM shelfs WHERE shelf_type = $1 LIMIT $2", item.MainShelf, limit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var arrIds pq.Int64Array
		err := rows.Scan(&shelf.Id, &shelf.ShelfType, &arrIds)
		if err != nil {
			fmt.Printf("error to scan rows: [%s]\n", err.Error())
		}
		shelf.Items = arrIds

		shelfs = append(shelfs, shelf)
	}

	// for _, shelf := range shelfs {
	// 	for _, i := range shelf.Items {
	// 		if i == item.Id {
	// 			if quantityMap[item.ItemName] == quantity {
	// 				break
	// 			}
	// 			quantityMap[item.ItemName] += 1
	// 		}
	// 	}

	// }

	return nil
}

func (p *Postgres) updateShelfState(shelfId int64, updatedItems []int64) error {
	_, err := p.db.Exec(`UPDATE shelfs SET items = $1 WHERE id = $2`, shelfId, updatedItems)
	if err != nil {
		return fmt.Errorf("error to update the shelf items state: [%s]", err.Error())
	}

	return nil
}

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

	psql := NewPostgres(db)

	// if err := psql.createOrders(); err != nil {
	// 	log.Println(err)
	// }

	// Get formed orders description from database
	order, err := psql.getOrder(args)
	if err != nil {
		log.Println(err)
	}

	// Get particular order description
	if err := psql.getOrderDescription(order); err != nil {
		log.Println(err)
	}

	// shelfs, err := psql.getShelfs()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// if err := psql.rangeOverOrders(order); err != nil {
	// 	log.Println(err)
	// }

	// if err := repo.createRecords(paths); err != nil {
	// 	fmt.Println(err)
	// }

	// if err := repo.getItems(res, args); err != nil {
	// 	fmt.Println(err)
	// }
}

// func (p *Postgres) getItemsFromOtherShelfs(shelfType string, limit int64) error {
// 	var shelf Shelf
// 	var shelfs []Shelf
// 	rows, err := p.db.Query("SELECT * FROM shelfs WHERE shelf_type = $1 LIMIT $2", shelfType, limit)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	return nil
// }

// if item.OtherShelfs != nil {
// 	for _, ty := range item.OtherShelfs {
// 		err := p.getItemsFromOtherShelfs(string(ty), 10)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// }

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

func (p *Postgres) selectItemById(itemId int64) (*Item, error) {
	var item Item
	err := p.db.QueryRow("SELECT * FROM items WHERE id = $1", itemId).
		Scan(&item.Id, &item.ItemName, &item.MainShelf, &item.OtherShelfs)
	if err != nil {
		return nil, fmt.Errorf("error to get item by id[%d]: [%s]", itemId, err.Error())
	}
	return &item, nil
}

func (p *Postgres) createShelfs() error {
	query := `CREATE TABLE IF NOT EXISTS shelfs(
		id BIGSERIAL PRIMARY KEY NOT NULL,
		shelf_type CHAR,
		items INTEGER[]);
	INSERT INTO shelfs(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
	INSERT INTO shelfs(shelf_type, items) VALUES('Б', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 3, 8, 8, 8, 8]);
	INSERT INTO shelfs(shelf_type, items) VALUES('Ж', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 6, 8, 8, 8, 4]);
	INSERT INTO shelfs(shelf_type, items) VALUES('З', ARRAY[3, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
	INSERT INTO shelfs(shelf_type, items) VALUES('В', ARRAY[3, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
	INSERT INTO shelfs(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
	INSERT INTO shelfs(shelf_type, items) VALUES('Б', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 3, 8, 8, 8, 8]);
	INSERT INTO shelfs(shelf_type, items) VALUES('Ж', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 6, 8, 8, 8, 4]);
	INSERT INTO shelfs(shelf_type, items) VALUES('З', ARRAY[3, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
	INSERT INTO shelfs(shelf_type, items) VALUES('В', ARRAY[3, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
	INSERT INTO shelfs(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
	INSERT INTO shelfs(shelf_type, items) VALUES('Б', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 3, 8, 8, 8, 8]);
	INSERT INTO shelfs(shelf_type, items) VALUES('Ж', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 6, 8, 8, 8, 4]);
	INSERT INTO shelfs(shelf_type, items) VALUES('З', ARRAY[3, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
	INSERT INTO shelfs(shelf_type, items) VALUES('В', ARRAY[3, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
	INSERT INTO shelfs(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);`
	if _, err := p.db.Exec(query); err != nil {
		return fmt.Errorf("error to insert shelfs: [%s]", err.Error())
	}
	return nil
}

func (p *Postgres) createOrders() error {
	query := `CREATE TABLE IF NOT EXISTS orders(
		id INTEGER,
		item_id INTEGER,
		quantity INTEGER);
	INSERT INTO orders(id, item_id, quantity) VALUES(10, 1, 2);
	INSERT INTO orders(id, item_id, quantity) VALUES(10, 3, 1);
	INSERT INTO orders(id, item_id, quantity) VALUES(10, 6, 1);
	INSERT INTO orders(id, item_id, quantity) VALUES(11, 2, 3);
	INSERT INTO orders(id, item_id, quantity) VALUES(14, 1, 3);
	INSERT INTO orders(id, item_id, quantity) VALUES(14, 2, 4);
	INSERT INTO orders(id, item_id, quantity) VALUES(15, 5, 1);`
	if _, err := p.db.Exec(query); err != nil {
		return fmt.Errorf("error to insert orders: [%s]", err.Error())
	}
	return nil
}

// func (p *Postgres) createRecords(paths []string) error {
// 	for _, path := range paths {
// 		file, err := os.ReadFile(path)
// 		if err != nil {
// 			return fmt.Errorf("error to read the file: [%s]", err.Error())
// 		}

// 		rows := strings.Split(string(file), ";")
// 		for _, row := range rows {
// 			_, err := p.db.Exec(row)
// 			if err != nil {
// 				fmt.Printf("error executing SQL statement - [%s], err: [%s]\n", row, err.Error())
// 			}
// 		}
// 	}
// 	return nil
// }

// func (p *Postgres) getItems(shelfs []Shelf, args []string) error {
// 	for _, arg := range args {
// 		argToInt, err := strconv.Atoi(arg)
// 		if err != nil {
// 			fmt.Printf("error to cast arg to integer: [%s]\n", err.Error())
// 		}
// 		for _, shelf := range shelfs {
// 			id := shelf.Items[argToInt-1]
// 			itemsMap := make(map[int64]int)

// 			for _, total := range shelf.Items {
// 				itemsMap[total] += 1
// 			}

// 			// count := itemsMap[id]

// 			item, err := p.selectItemById(id)
// 			if err != nil {
// 				fmt.Printf("error to get item by given id: [%s]\n", err.Error())
// 			}

// 			fmt.Println(item)
// 		}
// 	}

// 	return nil
// }
