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
				// var offset int64 = 0
				var limit int64 = 1

				quantity, err = s.repository.getItemsFromShelf(item, reserve.Quantity, limit)
				if err != nil {
					log.Println(err)
				}

				// if quantity != reserve.Quantity {

				// 	if len(item.OtherShelfs) == 0 {

				// 		for quantity != reserve.Quantity {
				// 			time.Sleep(1 * time.Second)
				// 			limit += 5
				// 			offset = limit

				// 			quantity, err = s.repository.getItemsFromShelf(item, reserve.Quantity, offset, limit, quantity)
				// 			if err != nil {
				// 				log.Println(err)
				// 			}
				// 		}

				// 	} else {
				// 		// quantity, err = s.repository.getItemsFromOtherShelfsWithLimit([]uint8{})
				// 		// if err != nil {
				// 		// 	log.Println(err)
				// 		// }
				// 		fmt.Println("got on two")
				// 	}
				// }

				itemOthShel := []string{}
				if item.OtherShelfs != nil {
					for _, shel := range item.OtherShelfs {
						itemOthShel = append(itemOthShel, string(shel))
					}
				}

				order := newOrder(item.ItemName, item.Id, reserve.Id, quantity, itemOthShel)
				orders[string(item.MainShelf)] = append(orders[string(item.MainShelf)], *order)

			}
		}
	}

	return orders, nil
}
