package main

type Reserve struct {
	Id       int64
	ItemId   int64
	Quantity int64
}

type Item struct {
	Id          int64
	ItemName    string
	MainShelf   []uint8
	OtherShelfs []uint8
}

type Shelf struct {
	Id        int64
	ShelfType []uint8
	Items     []int64
}

type Order struct {
	ItemName    string
	ItemId      int64
	OrderId     int64
	Quantity    int64
	OtherShelfs []string
}

func newOrder(itName string, itId, ordId, qt int64, othShel []string) *Order {
	return &Order{
		ItemName:    itName,
		ItemId:      itId,
		OrderId:     ordId,
		Quantity:    qt,
		OtherShelfs: othShel,
	}
}
