package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) execute(orderNumbers []string) (map[string][]Order, error) {
	reserves, err := s.repository.getReserves()
	if err != nil {
		return nil, err
	}

	sort.Slice(reserves, func(i, j int) bool {
		return reserves[i].ItemId < reserves[j].ItemId
	})

	orders := make(map[string][]Order)
	for _, orderNumber := range orderNumbers {
		ordNum, err := strconv.Atoi(orderNumber)
		if err != nil {
			fmt.Printf("error to convert arg to int type: [%s]\n", err.Error())
			continue
		}

		for _, reserve := range reserves {
			if ordNum == int(reserve.Id) {
				item, err := s.repository.getItemById(reserve.ItemId)
				if err != nil {
					log.Println(err)
				}

				var quantity int64 = 0
				var offset int64 = 0
				var limit int64 = 5

				quantity, err = s.repository.getItemsFromShelf(item, reserve.Quantity)
				if err != nil {
					log.Println(err)
				}

				// If we get not expected quantity.
				if quantity != reserve.Quantity {
					// If we have additional shelfs of item.
					if len(item.OtherShelfs) != 0 {
						// Range over them (maybe there is multiple additional shelfs).
						for _, otherShel := range item.OtherShelfs {
							// Try to find with additional shelfs, and some limit.
							quantity, err = s.repository.getItemsFromOtherItemShelfs(item, quantity, reserve.Quantity, offset, limit, otherShel)
							if err != nil {
								log.Println(err)
							}

							// If not found in additional, start loop on main shelfs to find items.
							if quantity != reserve.Quantity {
								for quantity != reserve.Quantity {
									quantity, err = s.repository.getItemsFromShelfWithLimit(item, quantity, reserve.Quantity, offset, limit)
									if err != nil {
										log.Println(err)
									}
									limit += 5
									offset += 5
								}
							}
						}
					}
					// If doesn't have some additional shelfs, go immediately in loop to find items.
					for quantity != reserve.Quantity {
						quantity, err = s.repository.getItemsFromShelfWithLimit(item, quantity, reserve.Quantity, offset, limit)
						if err != nil {
							log.Println(err)
						}
						limit += 5
						offset += 5
					}
				}

				itemOthShel := []int32{}
				if item.OtherShelfs != nil {
					itemOthShel = append(itemOthShel, item.OtherShelfs...)
				}

				order := newOrder(item.ItemName, item.Id, reserve.Id, quantity, itemOthShel)
				orders[string(item.MainShelf)] = append(orders[string(item.MainShelf)], *order)

			}
		}
	}

	return orders, nil
}
