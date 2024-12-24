package services

import (
	"ClubmineStoreService/models"
	"ClubmineStoreService/models/RequestTypes"
	"ClubmineStoreService/stores"
)

type (
	ItemService interface {
		GetItemFromId(id string) (models.Item, error)
		GetItemFromStack(stack string) (models.Item, error)
		CreateItem(item *RequestTypes.CreateItemRequest) error
		DeleteItem(id string) error
	}

	itemService struct {
		stores *stores.ItemStore
	}
)

func (i itemService) GetItemFromId(id string) (models.Item, error) {
	r, err := stores.ItemStore.GetItemFromId(*i.stores, id)
	return r, err
}

// GetItemFromStack does not work. character limit.
func (i itemService) GetItemFromStack(stack string) (models.Item, error) {
	r, err := stores.ItemStore.GetItemFromStack(*i.stores, stack)
	return r, err
}

func (i itemService) CreateItem(item *RequestTypes.CreateItemRequest) error {
	err := stores.ItemStore.CreateItem(*i.stores, item)
	return err
}

func (i itemService) DeleteItem(id string) error {
	err := stores.ItemStore.DeleteItem(*i.stores, id)
	return err
}
