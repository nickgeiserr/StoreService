package RequestTypes

type GetListingRequest struct {
	ListingID string `json:"listing_id"`
	Page      int    `json:"page"`
	ListerID  string `json:"lister_id"`
}
