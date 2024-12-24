package services

import (
	"ClubmineStoreService/models"
	"ClubmineStoreService/stores"
)

type (
	ListingsService interface {
		GetAllListings(pageNumber int) ([]models.ListingResponse, error)
		GetAllWithFilter(id string, pageNum int, column string) ([]models.ListingResponse, error)
		GetListing(listingId string) (models.ListingResponse, error)
		CreateListing(listing models.Listing) (bool, error)
		DeleteListing(listingId string) error
		GetItemFromID(itemId string) (models.Item, error)
	}

	listingsService struct {
		stores *stores.ListingStore
	}
)

func (l listingsService) GetAllListings(pageNumber int) ([]models.ListingResponse, error) {
	r, err := stores.ListingStore.GetAll(*l.stores, pageNumber)
	return r, err
}

func (l listingsService) GetAllWithFilter(id string, pageNum int, column string) ([]models.ListingResponse, error) {
	switch column {
	case "lister":
		r, err := stores.ListingStore.GetAllByListerID(*l.stores, id, pageNum)
		return r, err
	}

	return nil, nil
}

func (l listingsService) GetListing(listingId string) (models.ListingResponse, error) {
	r, err := stores.ListingStore.GetListing(*l.stores, listingId)
	return r, err
}

func (l listingsService) CreateListing(listing models.Listing) (bool, error) {
	r, err := stores.ListingStore.CreateListing(*l.stores, &listing)
	return r, err
}

func (l listingsService) DeleteListing(listingId string) error {
	err := stores.ListingStore.DeleteListing(*l.stores, listingId)
	return err
}

func (l listingsService) GetItemFromID(itemId string) (models.Item, error) {
	r, err := stores.ListingStore.GetItemFromID(*l.stores, itemId)
	return r, err
}
