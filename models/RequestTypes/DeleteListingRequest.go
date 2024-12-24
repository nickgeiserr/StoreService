package RequestTypes

type DeleteListingRequest struct {
	ID        string `json:"requester_id"`
	ListingID string `json:"listing_id"`
}
