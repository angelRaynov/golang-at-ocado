package main

import (
	"context"
	"errors"
	"github.com/angelRaynov/golang-at-ocado/week-03/fulfillment-skeleton/sort/gen"
	"google.golang.org/grpc"
	"log"
)

var Conn *grpc.ClientConn
var FConn *grpc.ClientConn

var loadItemPayload = &gen.LoadItemsRequest{
	Items: []*gen.Item{
		{Code: "123", Label: "tomato"},
		{Code: "456", Label: "cucumber"},
		{Code: "420", Label: "glass"},
		{Code: "222", Label: "fork"},
		{Code: "111", Label: "english breakfast"},
		{Code: "333", Label: "beans in a can"},
		{Code: "666", Label: "peaches"},
		{Code: "667", Label: "oranges"},

		{Code: "501", Label: "headphones"},
		{Code: "502", Label: "keyboard"},
		{Code: "503", Label: "cat in a box"},
		//6
		{Code: "601", Label: "book 1"},
		{Code: "602", Label: "book 2"},
		{Code: "603", Label: "book 3"},
		{Code: "604", Label: "book 4"},
		//7
		{Code: "401", Label: "water bottle"},
		{Code: "402", Label: "wataaa"},
		//8
		{Code: "301", Label: "juice"},
		//9
		{Code: "201", Label: "toy"},
		{Code: "202", Label: "teddy bear"},
		{Code: "203", Label: "dinosaur"},
		{Code: "204", Label: "dog"},
		{Code: "205", Label: "mug"},
		//10
		{Code: "101", Label: "laptop"},
		{Code: "102", Label: "mouse"},
	},
}
var loadOrderPayload = &gen.LoadOrdersRequest{
	Orders: []*gen.Order{
		{
			Id: "1", Items: []*gen.Item{
			{Code: "123", Label: "tomato"},
			{Code: "456", Label: "cucumber"},
		},
		},
		{
			Id: "2", Items: []*gen.Item{
			{Code: "420", Label: "glass"},
			{Code: "222", Label: "fork"},
		},
		},
		{
			Id: "3", Items: []*gen.Item{
			{Code: "111", Label: "english breakfast"},
			{Code: "333", Label: "beans in a can"},
		},
		},
		{
			Id: "4", Items: []*gen.Item{
			{Code: "666", Label: "peaches"},
			{Code: "667", Label: "oranges"},
		},
		},
		{
			Id: "5", Items: []*gen.Item{
			{Code: "501", Label: "headphones"},
			{Code: "502", Label: "keyboard"},
			{Code: "503", Label: "cat in a box"},
		},
		},
		{
			Id: "6", Items: []*gen.Item{
			{Code: "601", Label: "book 1"},
			{Code: "602", Label: "book 2"},
			{Code: "603", Label: "book 3"},
			{Code: "604", Label: "book 4"},
		},
		},
		{
			Id: "7", Items: []*gen.Item{
			{Code: "401", Label: "water bottle"},
			{Code: "402", Label: "wataaa"},
		},
		},
		{
			Id: "8", Items: []*gen.Item{
			{Code: "301", Label: "juice"},
		},
		},
		{
			Id: "9", Items: []*gen.Item{
			{Code: "201", Label: "toy"},
			{Code: "202", Label: "teddy bear"},
			{Code: "203", Label: "dinosaur"},
			{Code: "204", Label: "dog"},
			{Code: "205", Label: "mug"},
		},
		},
		{
			Id: "10", Items: []*gen.Item{
			{Code: "101", Label: "laptop"},
			{Code: "102", Label: "mouse"},
		},
		},
	},
}

func main() {
	initRobotConnection()
	defer Conn.Close()

	robot := gen.NewSortingRobotClient(Conn)

	initFulfillmentConnection()
	defer FConn.Close()

	fulfillment := gen.NewFulfillmentClient(FConn)

	ctx := context.Background()

	//load items
	_, err := robot.LoadItems(ctx, loadItemPayload)
	if err != nil {
		log.Fatal(err)
	}

	//load orders
	orderResponse, err := fulfillment.LoadOrders(ctx, loadOrderPayload)
	if err != nil {
		log.Fatal(err)
	}

	for {

		//pick item
		item, err := robot.PickItem(ctx, &gen.Empty{})
		if err != nil {
			log.Print(err)
			break
		}

		pickedItem := item.Item.Label
		log.Printf("Picked item %#v", item.Item.Label)

		//match
		prepOrd, err := matchOrder(pickedItem, orderResponse)
		if err != nil {
			log.Print(err)
			return
		}

		//dispatch
		_, err = robot.PlaceInCubby(ctx, getMoveItemPayload(prepOrd.Cubby.Id))
		if err != nil {
			log.Print(err)
			return
		}
	}

}


func getMoveItemPayload(id string) *gen.PlaceInCubbyRequest {
	return &gen.PlaceInCubbyRequest{
		Cubby: &gen.Cubby{Id: id},
	}
}

func initRobotConnection() {
	var err error
	Conn, err = grpc.Dial(":10000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect %s", err)
	}

}

func initFulfillmentConnection() {
	var err error
	FConn, err = grpc.Dial(":10001", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect %s", err)
	}

}

func matchOrder(pickedItem string, resp *gen.CompleteResponse) (*gen.PreparedOrder, error) {
	for _, prepOrd := range resp.Orders {
		for _, currItem := range prepOrd.Order.Items {
			if pickedItem == currItem.Label {
				log.Printf("order %v -> cubby %v", prepOrd.Order.Id, prepOrd.Cubby.Id)
				return prepOrd, nil
			}
		}

	}
	return nil, errors.New("no match for the picked item")
}
