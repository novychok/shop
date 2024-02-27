package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/lib/pq"
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
	Shelf    string
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

	// sort.Slice(order.OrderDescription, func(i, j int) bool {
	// 	return order.OrderDescription[i].Id < order.OrderDescription[j].Id
	// })

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

func (p *Postgres) getOrderDescription(order *Order) ([]FinalOrder, error) {
	var finalOrders []FinalOrder
	for _, orderDesc := range order.OrderDescription {

		item, err := p.selectItemById(orderDesc.ItemId)
		if err != nil {
			log.Println(err)
		}

		// Get products from shelfs
		quantity, err := p.getItemsFromShelf(item)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(quantity)

		// Create final order
		finalOrder := FinalOrder{
			Shelf:    string(item.MainShelf),
			Products: []Product{},
		}
		product := Product{
			ItemName:    item.ItemName,
			ItemId:      item.Id,
			OrderNumber: orderDesc.Id,
			ItemCount:   quantity,
		}

		// Start to search in other shelfs if we don't get enough quantity of products
		if quantity != orderDesc.Quantity {

			if item.OtherShelfs != nil {
				// Take items from item other shelfs
				quantity, err := p.getItemsFromOtherShelf(item, quantity, orderDesc.Quantity, 100)
				if err != nil {
					log.Println(err)
				}

				if quantity != orderDesc.Quantity {
					// Make recursive search of items from shelfs
					fmt.Println("sorry don't have enough items in stock")
					return finalOrders, nil
				}

			}

		}

		finalOrder.Products = append(finalOrder.Products, product)
		finalOrders = append(finalOrders, finalOrder)
	}

	return finalOrders, nil
}

// Get items from item other shelfs, if we don't get enough item quantity from shelf
// and if item other shelf != nil
func (p *Postgres) getItemsFromOtherShelf(item *Item, count, orderCount, limit int64) (int64, error) {
	var shelf Shelf
	var arrIds pq.Int64Array
	var quantity int64 = count
	for _, ty := range item.OtherShelfs {
		rows, err := p.db.Query("SELECT * FROM shelfs WHERE shelf_type = $1 LIMIT $2", ty, limit)
		if err != nil {
			return 0, fmt.Errorf("error getItemsFromOtherShelf to scan row: [%s]", err.Error())
		}
		defer rows.Close()
		shelf.Items = arrIds // Use arrIds instead this

		for rows.Next() {
			err := rows.Scan(&shelf.Id, &shelf.ShelfType, &arrIds)
			if err != nil {
				log.Fatal(err)
			}

			for i, id := range shelf.Items {
				if id != item.Id {
					continue
				}
				quantity += 1
				shelf.Items[i] = 0
				if err := p.updateShelfState(shelf.Id, shelf.Items); err != nil {
					log.Println(err)
				}

				if quantity == orderCount {
					return quantity, nil
				}
			}
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

	}

	return quantity, nil
}

// Get items from a shelf
func (p *Postgres) getItemsFromShelf(item *Item) (int64, error) {
	var shelf Shelf
	var arrIds pq.Int64Array
	err := p.db.QueryRow("SELECT * FROM shelfs WHERE shelf_type = $1 LIMIT 1", item.MainShelf).
		Scan(&shelf.Id, &shelf.ShelfType, &arrIds)
	if err != nil {
		return 0, fmt.Errorf("error getItemsFromShelf to scan row: [%s]", err.Error())
	}
	shelf.Items = arrIds

	var quantity int64 = 0
	for i, id := range shelf.Items {
		if id != item.Id {
			continue
		}
		quantity += 1
		shelf.Items[i] = 0
		if err := p.updateShelfState(shelf.Id, shelf.Items); err != nil {
			log.Println(err)
		}
	}

	return quantity, nil
}

// Update shelf items after retrieve certain item
func (p *Postgres) updateShelfState(shelfId int64, updatedItems []int64) error {
	var ret pq.Int64Array
	ret = updatedItems
	_, err := p.db.Exec(`UPDATE shelfs SET items = $1 WHERE id = $2`, ret, shelfId)
	if err != nil {
		return fmt.Errorf("error to update the shelf items state: [%s]", err.Error())
	}

	return nil
}

// Get single item from psql
func (p *Postgres) selectItemById(itemId int64) (*Item, error) {
	var item Item
	err := p.db.QueryRow("SELECT * FROM items WHERE id = $1", itemId).
		Scan(&item.Id, &item.ItemName, &item.MainShelf, &item.OtherShelfs)
	if err != nil {
		return nil, fmt.Errorf("error to get item by id[%d]: [%s]", itemId, err.Error())
	}
	return &item, nil
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

	// Init tables orders, shelfs, items
	if err := psql.initTables(); err != nil {
		log.Println(err)
	}

	order, err := psql.getOrder(args)
	if err != nil {
		log.Println(err)
	}

	finalOrders, err := psql.getOrderDescription(order)
	if err != nil {
		log.Println(err)
	}

	for _, finalOrder := range finalOrders {
		fmt.Printf("%+v\n", finalOrder)
	}
}

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

func (p *Postgres) initTables() error {
	if err := p.createOrders(); err != nil {
		log.Println(err)
	}
	if err := p.createShelfs(); err != nil {
		log.Println(err)
	}
	if err := p.createItems(); err != nil {
		log.Println(err)
	}

	return nil
}

func (p *Postgres) createOrders() error {
	query := `CREATE TABLE IF NOT EXISTS order_description(
		id INTEGER,
		item_id INTEGER,
		quantity INTEGER);
	INSERT INTO order_description(id, item_id, quantity) VALUES(10, 1, 2);
	INSERT INTO order_description(id, item_id, quantity) VALUES(10, 3, 1);
	INSERT INTO order_description(id, item_id, quantity) VALUES(10, 6, 1);
	INSERT INTO order_description(id, item_id, quantity) VALUES(11, 2, 3);
	INSERT INTO order_description(id, item_id, quantity) VALUES(14, 1, 3);
	INSERT INTO order_description(id, item_id, quantity) VALUES(14, 2, 4);
	INSERT INTO order_description(id, item_id, quantity) VALUES(15, 5, 1);`
	if _, err := p.db.Exec(query); err != nil {
		return fmt.Errorf("error to insert orders: [%s]", err.Error())
	}
	return nil
}

func (p *Postgres) createShelfs() error {
	query := `CREATE TABLE IF NOT EXISTS shelfs(
		id BIGSERIAL PRIMARY KEY NOT NULL,
		shelf_type CHAR,
		items INTEGER[]);
	INSERT INTO shelfs(shelf_type, items) VALUES('А', ARRAY[1,1,3,6,2,2,2,1,1,1,2,2,2,2,5]);
	INSERT INTO shelfs(shelf_type, items) VALUES('А', ARRAY[1,1,3,6,2,2,2,1,1,1,2,2,2,2,5]);
	INSERT INTO shelfs(shelf_type, items) VALUES('Б', ARRAY[1,1,3,6,2,2,2,1,1,1,2,2,2,2,5]);
	INSERT INTO shelfs(shelf_type, items) VALUES('Б', ARRAY[1,1,3,6,2,2,2,1,1,1,2,2,2,2,5]);
	INSERT INTO shelfs(shelf_type, items) VALUES('Ж', ARRAY[1,1,3,6,2,2,2,1,1,1,2,2,2,2,5]);
	INSERT INTO shelfs(shelf_type, items) VALUES('Ж', ARRAY[1,1,3,6,2,2,2,1,1,1,2,2,2,2,5]);
	INSERT INTO shelfs(shelf_type, items) VALUES('З', ARRAY[1,1,3,6,2,2,2,1,1,1,2,2,2,2,5]);
	INSERT INTO shelfs(shelf_type, items) VALUES('З', ARRAY[1,1,3,6,2,2,2,1,1,1,2,2,2,2,5]);
	INSERT INTO shelfs(shelf_type, items) VALUES('В', ARRAY[1,1,3,6,2,2,2,1,1,1,2,2,2,2,5]);
	INSERT INTO shelfs(shelf_type, items) VALUES('В', ARRAY[1,1,3,6,2,2,2,1,1,1,2,2,2,2,5]);`
	if _, err := p.db.Exec(query); err != nil {
		return fmt.Errorf("error to insert shelfs: [%s]", err.Error())
	}
	return nil
}

func (p *Postgres) createItems() error {
	query := `CREATE TABLE IF NOT EXISTS items(
		id BIGSERIAL PRIMARY KEY NOT NULL,
		item_name VARCHAR(100),
		main_shelf CHAR,
		other_shelfs CHAR[]);
	INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Laptop', 'А', NULL);
	INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Moitor', 'А', NULL);
	INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Phone', 'Б', '{З,В}');
	INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('PC', 'Ж', NULL);
	INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Watch', 'Ж', NULL);
	INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Microphone', 'Ж', NULL);`
	if _, err := p.db.Exec(query); err != nil {
		return fmt.Errorf("error to insert items: [%s]", err.Error())
	}
	return nil
}
