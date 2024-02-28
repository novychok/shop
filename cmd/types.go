package main

import "github.com/lib/pq"

type Reserve struct {
	Id       int64
	ItemId   int64
	Quantity int64
}

type Item struct {
	Id          int64
	ItemName    string
	MainShelf   []uint8
	OtherShelfs []rune
}

type Shelf struct {
	Id        int64
	ShelfType []uint8
	Items     pq.Int64Array
}

type Order struct {
	ItemName    string
	ItemId      int64
	OrderId     int64
	Quantity    int64
	OtherShelfs []rune
}

func newOrder(itName string, itId, ordId, qt int64, othShel []rune) *Order {
	return &Order{
		ItemName:    itName,
		ItemId:      itId,
		OrderId:     ordId,
		Quantity:    qt,
		OtherShelfs: othShel,
	}
}
