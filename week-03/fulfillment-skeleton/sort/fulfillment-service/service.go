package main

import (
	"context"
	"errors"
	"github.com/preslavmihaylov/ordertocubby"
	"log"
	"math/rand"
	"time"

	"github.com/angelRaynov/golang-at-ocado/week-03/fulfillment-skeleton/sort/gen"
)

func newFulfillmentService(client gen.SortingRobotClient) gen.FulfillmentServer {
	rand.Seed(time.Now().UnixNano())
	return &fulfillmentService{
		sortingRobot: client,
	}
}

type fulfillmentService struct {
	sortingRobot gen.SortingRobotClient
}

func (fs *fulfillmentService) LoadOrders(ctx context.Context, req *gen.LoadOrdersRequest) (*gen.CompleteResponse, error) {
	//ordersToCubbies := mapOrdersToCubbies(req.Orders)
	//log.Printf("%v\n", ordersToCubbies)
	//
	//for _, order := range req.Orders {
	//	for _, item := range order.Items {
	//		_ = item
	//
	//		pickedItem, err := fs.sortingRobot.PickItem(ctx, &gen.Empty{})
	//		if err != nil {
	//			return nil, errors.New("not implemented")
	//		}
	//
	//		cubbyID := getCubbyForItem(pickedItem.Item)
	//		log.Printf("picked item %v\n", pickedItem)
	//
	//		_, err = fs.sortingRobot.PlaceInCubby(ctx, &gen.PlaceInCubbyRequest{
	//			Cubby: &gen.Cubby{Id: cubbyID},
	//		})
	//
	//		if err != nil {
	//			return nil, fmt.Errorf("place in cubby failed: %v", err)
	//		}
	//	}
	//}

	var bookedCubbies = make(map[string]bool)
	loadedOrders := getOrders(bookedCubbies, req)

	for {
		pickedItem, err := fs.sortingRobot.PickItem(ctx, &gen.Empty{})
		if err != nil {
			return nil, err
		}

		//match
		preparedOrder, err := matchOrder(pickedItem.Item.Label, loadedOrders)
		if err != nil {
			return nil, err
		}

		//dispatch
		_, err = fs.sortingRobot.PlaceInCubby(ctx, getMoveItemPayload(preparedOrder.Cubby.Id))
		if err != nil {
			return nil, err
		}
	}
}

//func mapOrdersToCubbies(orders []*gen.Order) map[string]string {
//	ordersToCubbies := map[string]string{}
//	usedCubbies := map[string]bool{}
//
//	for _, order := range orders {
//		cubbyID := getUniqueCubby(usedCubbies, order.Id, cubbiesCnt)
//		ordersToCubbies[order.Id] = cubbyID
//		usedCubbies[cubbyID] = true
//	}
//
//	for orderID, cubbyID := range ordersToCubbies {
//		fmt.Printf("order %s -> cubby %s\n", orderID, cubbyID)
//	}
//
//	return ordersToCubbies
//}
//
//func getUniqueCubby(usedCubbies map[string]bool, id string, cubbiesCnt int) string {
//	times := 1
//	for {
//		cubbyID := ordertocubby.Map(id, uint32(times), uint32(cubbiesCnt))
//		if !usedCubbies[cubbyID] {
//			return cubbyID
//		}
//
//		times++
//	}
//}
//
//func getCubbyForItem(item *gen.Item) string {
//	return "1"
//}

func getUniqueCubby(orderID string, times uint32, bookedCubbies map[string]bool) string {
	cubbyId := ordertocubby.Map(orderID, times, 10)

	if _, ok := bookedCubbies[cubbyId]; ok {
		return getUniqueCubby(orderID, times+1, bookedCubbies)
	}

	bookedCubbies[cubbyId] = true
	return cubbyId
}

func getOrders(bookedCubbies map[string]bool, req *gen.LoadOrdersRequest) []*gen.PreparedOrder {
	var orders []*gen.PreparedOrder

	for _, ord := range req.Orders {
		cubbyId := getUniqueCubby(ord.Id, 1, bookedCubbies)
		preparedOrder := &gen.PreparedOrder{Order: ord, Cubby: &gen.Cubby{Id: cubbyId}}
		orders = append(orders, preparedOrder)
	}

	return orders
}

func matchOrder(pickedItem string, loadedOrders []*gen.PreparedOrder) (*gen.PreparedOrder, error) {
	for _, preparedOrder := range loadedOrders {
		for _, currItem := range preparedOrder.Order.Items {
			if pickedItem == currItem.Label {
				log.Printf("order %v -> cubby %v", preparedOrder.Order.Id, preparedOrder.Cubby.Id)
				return preparedOrder, nil
			}
		}

	}
	return nil, errors.New("no match for the picked item")
}

func getMoveItemPayload(id string) *gen.PlaceInCubbyRequest {
	return &gen.PlaceInCubbyRequest{
		Cubby: &gen.Cubby{Id: id},
	}
}
