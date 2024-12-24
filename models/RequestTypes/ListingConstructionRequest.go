package RequestTypes

type ListingConstructionRequest struct {
	ListerID  string `json:"lister_id"`
	ItemID    string `json:"item_id"`
	ItemStack string `json:"item_stack"`
	Price     int    `json:"price"`
}
