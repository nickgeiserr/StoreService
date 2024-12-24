package models

type ListingResponse struct {
	ListingID string `json:"listing_id"`
	ListerID  string `json:"lister_id"`
	ItemID    string `json:"item_id"`
	ItemStack string `json:"item_stack"`
	Price     int    `json:"price"`
	CreatedAt string `json:"created_at"`
}
