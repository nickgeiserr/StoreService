package RequestTypes

type CreateItemRequest struct {
	ItemID    string `json:"item_id"`
	ItemStack string `json:"item_stack"`
	CreatorID string `json:"creator_id"`
}
