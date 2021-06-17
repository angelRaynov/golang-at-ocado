package main

import (
	"context"
	"github.com/angelRaynov/golang-at-ocado/week-03/fulfillment-skeleton/sort/gen"
	"google.golang.org/grpc"
	"log"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":10001", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect %s", err)
	}

	defer conn.Close()

	f := gen.NewFulfillmentClient(conn)

	payload := &gen.LoadOrdersRequest{
		Orders: []*gen.Order{
			{Id: "13",
				Items: []*gen.Item{
					{Code: "1", Label: "pecorino"},
					{Code: "2", Label: "grana padano"},
				},
			},
		},
	}

	ctx := context.Background()

	_, err = f.LoadOrders(ctx, payload)
	if err != nil {
		log.Print(err)
	}

}
