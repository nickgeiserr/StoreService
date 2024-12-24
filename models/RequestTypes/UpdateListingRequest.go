package RequestTypes

type UpdateListingRequest struct {
	ListingID string `json:"listing_id"`
	ID        string `json:"requester_id"`
	Price     int    `json:"price"`
}
