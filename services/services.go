package services

import "ClubmineStoreService/stores"

type Services struct {
	// User    UserService
	Listing  ListingsService
	Purchase PurchaseService
	Item     ItemService
}

func New(s *stores.Stores) *Services {
	return &Services{
		Listing:  &listingsService{&s.Listing},
		Purchase: &purchaseService{&s.Purchase},
		Item:     &itemService{&s.Item},
	}
}
