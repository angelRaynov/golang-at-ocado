package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/angelRaynov/golang-at-ocado/week-03/fulfillment-skeleton/sort/gen"
	"github.com/preslavmihaylov/ordertocubby"
)

func newFulfillmentService() gen.FulfillmentServer {
	rand.Seed(time.Now().UnixNano())
	return &fulfillmentService{}
}

type fulfillmentService struct {
}

func (s *fulfillmentService) LoadOrders(ctx context.Context, req *gen.LoadOrdersRequest) (*gen.CompleteResponse, error) {
	var bookedCubbies = make(map[string]interface{})
	var orders []*gen.PreparedOrder

	for _, ord := range req.Orders {
		cubbyId := getUniqueCubby(ord.Id, 1, bookedCubbies)
		preparedOrder := &gen.PreparedOrder{Order: ord, Cubby: &gen.Cubby{Id: cubbyId}}
		orders = append(orders,preparedOrder)
	}

	return &gen.CompleteResponse{Orders: orders, Status: "unfulfilled"}, nil
}

func getUniqueCubby(orderID string, times uint32, bookedCubbies map[string]interface{}) string {
	cubbyId := ordertocubby.Map(orderID, times, 10)

	if _, ok := bookedCubbies[cubbyId]; ok {
		return getUniqueCubby(orderID, times+1, bookedCubbies)
	}

	bookedCubbies[cubbyId] = nil
	return cubbyId
}
