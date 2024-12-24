package models

import "time"

type Transaction struct {
	ListingID string    `json:"listing_id"`
	Lister    string    `json:"lister_id"`
	Buyer     string    `json:"buyer_id"`
	Price     int       `json:"price"`
	Date      time.Time `json:"date"`
}
