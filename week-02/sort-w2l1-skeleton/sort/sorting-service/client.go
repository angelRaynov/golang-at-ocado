package main

import (
	"context"
	"fmt"
	"github.com/angelRaynov/golang-at-ocado/week-02/sort-w2l1-skeleton/sort/gen"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"time"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":10000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect %s", err)
	}

	defer conn.Close()

	robot := gen.NewSortingRobotClient(conn)

	loadItemPayload := &gen.LoadItemsRequest{
		Items: []*gen.Item{
			{Code: "1", Label: "pecorino"},
			{Code: "2", Label: "grana padano"},
			{Code: "3", Label: "parmigiano reggiano"},
			{Code: "4", Label: "mozarella"},
			{Code: "5", Label: "scamorza"},
		},
	}

	ctx := context.Background()

	_, err = robot.LoadItems(ctx, loadItemPayload)
	if err != nil {
		log.Print(err)
	}

	for {

		_, err = robot.SelectItem(ctx, &gen.SelectItemRequest{})
		if err != nil {
			log.Print(err)
			break
		}

		rand.Seed(time.Now().UnixNano())
		cubbyId := fmt.Sprintf("%d", rand.Intn(10))

		movePayload := &gen.MoveItemRequest{
			Cubby: &gen.Cubby{Id: cubbyId},
		}

		_, err = robot.MoveItem(ctx, movePayload)
		if err != nil {
			log.Print(err)
			break
		}

	}

}
