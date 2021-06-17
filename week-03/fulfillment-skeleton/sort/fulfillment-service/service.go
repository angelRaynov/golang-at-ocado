package main

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/angelRaynov/golang-at-ocado/week-03/fulfillment-skeleton/sort/gen"
)

func newFulfillmentService() gen.FulfillmentServer {
	rand.Seed(time.Now().UnixNano())
	return &fulfillmentService{}
}

type fulfillmentService struct {
}

func (s *fulfillmentService) LoadOrders(context.Context, *gen.LoadOrdersRequest) (*gen.CompleteResponse, error) {
	return nil, errors.New("not implemented")
}
