package main

import (
	"database/sql"
	"fmt"
	"log"
	"unsafe"

	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// getItemsFromOtherItemShelfs make same functionality as getItemsFromShelfWithLimit, but use other shelfs, instead using main shelf.
// and have fixed limit (maybe there is multiple other shelfs)
func (r *Repository) getItemsFromOtherItemShelfs(item *Item, howMuchHave, howMuchNeed, offset, limit int64, otherShel int32) (int64, error) {
	otherShelToBytes := *(*[]uint8)(unsafe.Pointer(&otherShel)) // Casting to DB type to find with additional shelf
	fmt.Println(otherShelToBytes)
	rows, err := r.db.Query("SELECT * FROM shelfs WHERE shelf_type = $1 OFFSET $2 LIMIT $3", []uint8(otherShelToBytes), offset, limit)
	if err != nil {
		return 0, fmt.Errorf("error to execute the query: [%s]", err.Error())
	}

	for rows.Next() {
		var shelf Shelf
		err := rows.Scan(&shelf.Id, &shelf.ShelfType, &shelf.Items)
		if err != nil {
			return 0, fmt.Errorf("error to scan row with limit: [%s]", err.Error())
		}
		fmt.Println(shelf)
		fmt.Println(howMuchHave)

		for i, id := range shelf.Items {
			if id != item.Id {
				continue
			}

			if howMuchHave == howMuchNeed {
				return howMuchHave, nil
			}
			howMuchHave += 1

			shelf.Items[i] = 0
			if err := r.updateShelfById(shelf.Id, shelf.Items); err != nil {
				log.Println(err)
			}
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return howMuchHave, nil
}

func (r *Repository) getItemsFromShelfWithLimit(item *Item, howMuchHave, howMuchNeed, offset, limit int64) (int64, error) {
	rows, err := r.db.Query("SELECT * FROM shelfs WHERE shelf_type = $1 OFFSET $2 LIMIT $3", item.MainShelf, offset, limit)
	if err != nil {
		return 0, fmt.Errorf("error to execute the query: [%s]", err.Error())
	}

	for rows.Next() {
		var shelf Shelf
		err := rows.Scan(&shelf.Id, &shelf.ShelfType, &shelf.Items)
		if err != nil {
			return 0, fmt.Errorf("error to scan row with limit: [%s]", err.Error())
		}

		for i, id := range shelf.Items {
			if id != item.Id {
				continue
			}

			if howMuchHave == howMuchNeed {
				return howMuchHave, nil
			}
			howMuchHave += 1

			shelf.Items[i] = 0
			if err := r.updateShelfById(shelf.Id, shelf.Items); err != nil {
				log.Println(err)
			}
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return howMuchHave, nil
}

func (r *Repository) getItemsFromShelf(item *Item, reservedQuantity int64) (quantity int64, err error) {
	var shelf Shelf
	var arrIds pq.Int64Array
	err = r.db.QueryRow("SELECT * FROM shelfs WHERE shelf_type = $1 LIMIT 1", item.MainShelf).
		Scan(&shelf.Id, &shelf.ShelfType, &arrIds)
	if err != nil {
		return 0, fmt.Errorf("error getItemsFromShelf to scan row: [%s]", err.Error())
	}
	shelf.Items = arrIds

	for i, id := range arrIds {
		if id != item.Id {
			continue
		}
		if quantity == reservedQuantity {
			return quantity, nil
		}
		quantity += 1

		shelf.Items[i] = 0
		if err := r.updateShelfById(shelf.Id, shelf.Items); err != nil {
			log.Println(err)
		}
	}

	return quantity, nil
}

func (r *Repository) getReserves() ([]Reserve, error) {
	rows, err := r.db.Query("SELECT * FROM reserves")
	if err != nil {
		return nil, fmt.Errorf("error whlie get all reserves: [%s]", err.Error())
	}
	defer rows.Close()

	var reserve Reserve
	var reserves []Reserve

	for rows.Next() {
		err := rows.Scan(&reserve.Id, &reserve.ItemId, &reserve.Quantity)
		if err != nil {
			fmt.Printf("error to scan rows: [%s]\n", err.Error())
		}

		reserves = append(reserves, reserve)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("rows.Err: [%s]\n", err.Error())
	}

	return reserves, nil
}

func (r *Repository) getItemById(itemId int64) (*Item, error) {
	var item Item
	var otherShel pq.Int32Array
	err := r.db.QueryRow("SELECT * FROM items WHERE id = $1", itemId).
		Scan(&item.Id, &item.ItemName, &item.MainShelf, &otherShel)
	if err != nil {
		return nil, fmt.Errorf("error to get item by id[%d]: [%s]", itemId, err.Error())
	}
	item.OtherShelfs = otherShel

	return &item, nil
}

func (r *Repository) updateShelfById(shelfId int64, updatedItems []int64) error {
	var update pq.Int64Array = updatedItems
	_, err := r.db.Exec(`UPDATE shelfs SET items = $1 WHERE id = $2`, update, shelfId)
	if err != nil {
		return fmt.Errorf("error to update the shelf items state: [%s]", err.Error())
	}

	return nil
}
