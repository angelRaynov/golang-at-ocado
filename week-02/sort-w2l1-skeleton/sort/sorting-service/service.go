package main

import (
	"context"
	"errors"
	"github.com/angelRaynov/golang-at-ocado/week-02/sort-w2l1-skeleton/sort/gen"
	"log"
	"math/rand"
	"time"
)

func newSortingService() gen.SortingRobotServer {
	return &sortingService{}
}

type sortingService struct {
	Selection *gen.Item
	Items     []*gen.Item
	Cubby     *gen.Cubby
}

func (s *sortingService) LoadItems(ctx context.Context, request *gen.LoadItemsRequest) (*gen.LoadItemsResponse, error) {
	s.Items = append(s.Items, request.Items...)

	log.Printf("Loaded items: %v", s.Items)
	return &gen.LoadItemsResponse{}, nil
}

func (s *sortingService) SelectItem(ctx context.Context, req *gen.SelectItemRequest) (*gen.SelectItemResponse, error) {
	if s.Selection != nil {
		return nil, errors.New("put that item in a cubby before picking another one")
	}

	if len(s.Items) == 0 {
		return nil, errors.New("can not select item from an empty bin")
	}

	rand.Seed(time.Now().Unix())
	idx := rand.Intn(len(s.Items))
	s.Selection = s.Items[idx]
	s.Items = append(s.Items[:idx], s.Items[idx+1:]...)

	log.Printf("Selected item %s", s.Selection.Label)

	if len(s.Items) != 0 {
		log.Printf("Remaining items in the bin %v", s.Items)
	} else {
		log.Println("The bin is empty")
	}

	return &gen.SelectItemResponse{}, nil
}

func (s *sortingService) MoveItem(ctx context.Context, req *gen.MoveItemRequest) (*gen.MoveItemResponse, error) {
	if req.Cubby == nil {
		s.returnItemToBin()
		return nil, errors.New("cubby not selected, returning the item")
	}

	if s.Selection == nil {
		return nil, errors.New("item not selected")
	}

	log.Printf("Item %s placed in cubby with id = %s", s.Selection.Label, req.Cubby.Id)
	s.Selection = nil

	return &gen.MoveItemResponse{}, nil
}

func (s *sortingService) returnItemToBin() {
	s.Items = append(s.Items, s.Selection)
}
