package main

import (
	"context"
	"github.com/angelRaynov/golang-at-ocado/week-02/sort-w2l1-skeleton/sort/gen"
	"testing"
)

func TestSortingService_LoadItems(t *testing.T) {
	s := &sortingService{}
	ctx := context.Background()

	payload := getLoadItemsPayload()

	_, err := s.LoadItems(ctx, payload)
	if err != nil {
		t.Error(err)
	}

	if len(s.Items) == 0 {
		t.Errorf("Items were loaded but the bin is empty")
	}
}

func TestSortingService_SelectItem(t *testing.T) {
	s := &sortingService{}
	ctx := context.Background()

	payload := getLoadItemsPayload()

	_, err := s.LoadItems(ctx, payload)
	if err != nil {
		t.Errorf("Unable to load items. %v", err)
	}

	expectedCount := len(s.Items) - 1
	_, err = s.SelectItem(ctx, &gen.SelectItemRequest{})
	if err != nil {
		t.Error(err)
	}

	count := len(s.Items)
	if expectedCount != count {
		t.Errorf("Expected count %d, got %d", expectedCount, count)
	}

}

func TestSortingService_LoadItemsDoubleSelect(t *testing.T) {
	s := &sortingService{}
	ctx := context.Background()

	payload := getLoadItemsPayload()

	_, err := s.LoadItems(ctx, payload)
	if err != nil {
		t.Errorf("Unable to load items. %v", err)
	}

	_, err = s.SelectItem(ctx, &gen.SelectItemRequest{})
	if err != nil {
		t.Error(err)
	}

	_, err = s.SelectItem(ctx, &gen.SelectItemRequest{})
	if err == nil {
		t.Errorf("Expected to put that item in a cubby before picking another one")
	}
}

func TestSortingService_SelectItemEmptyBin(t *testing.T) {
	s := &sortingService{}
	ctx := context.Background()

	_, err := s.SelectItem(ctx, &gen.SelectItemRequest{})
	if err == nil {
		t.Errorf("Expected empty bin error")
	}
}

func TestSortingService_MoveItem(t *testing.T) {
	s := &sortingService{}
	ctx := context.Background()

	payload := getLoadItemsPayload()

	_, err := s.LoadItems(ctx, payload)
	if err != nil {
		t.Errorf("Unable to load items. %v", err)
	}

	_, err = s.SelectItem(ctx, &gen.SelectItemRequest{})
	if err != nil {
		t.Errorf("Unable to select item. %v", err)
	}

	movePayload := &gen.MoveItemRequest{
		Cubby: &gen.Cubby{Id: "1"},
	}

	_, err = s.MoveItem(ctx, movePayload)
	if err != nil {
		t.Error(err)
	}

}

func TestSortingService_MoveItemNoCubby(t *testing.T) {
	s := &sortingService{}
	ctx := context.Background()

	payload := getLoadItemsPayload()

	_, err := s.LoadItems(ctx, payload)
	if err != nil {
		t.Errorf("Unable to load items. %v", err)
	}

	//no cubby, return the item
	noCubby := &gen.MoveItemRequest{
		Cubby: nil,
	}

	countBeforeSelect := len(s.Items)

	_, err = s.SelectItem(ctx, &gen.SelectItemRequest{})
	if err != nil {
		t.Errorf("Unable to select item. %v", err)
	}

	_, err = s.MoveItem(ctx, noCubby)
	if err == nil {
		t.Errorf("Cubby not selected, returning the item")
	}

	currentCount := len(s.Items)

	if currentCount != countBeforeSelect {
		t.Errorf("Item not returned to the bin")
	}
}

func TestSortingService_MoveItemNoSelectedItem(t *testing.T) {
	s := &sortingService{}
	ctx := context.Background()

	movePayload := &gen.MoveItemRequest{
		Cubby: &gen.Cubby{Id: "1"},
	}

	//no selected item
	s.Selection = nil
	_, err := s.MoveItem(ctx, movePayload)
	if err == nil {
		t.Errorf("Item not selected")
	}
}

func getLoadItemsPayload() *gen.LoadItemsRequest {
	return &gen.LoadItemsRequest{
		Items: []*gen.Item{
			{Code: "1", Label: "pecorino"},
			{Code: "2", Label: "grana padano"},
			{Code: "3", Label: "parmigiano reggiano"},
		},
	}
}
