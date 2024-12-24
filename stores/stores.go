package stores

type Stores struct {
	// User    UserStore
	Listing  ListingStore
	Purchase PurchaseStore
	Item     ItemStore
}

func New() *Stores {
	return &Stores{
		Listing:  &listingStore{},
		Purchase: &purchaseStore{},
		Item:     &itemStore{},
	}
}
