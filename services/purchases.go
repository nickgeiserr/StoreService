package services

import (
	"ClubmineStoreService/models"
	"ClubmineStoreService/stores"
	"errors"
)

type (
	PurchaseService interface {
		Purchase(l models.ListingResponse, transaction *models.TransactionAttempt) (bool, error)
	}

	purchaseService struct {
		stores *stores.PurchaseStore
	}
)

func (p purchaseService) Purchase(l models.ListingResponse, t *models.TransactionAttempt) (bool, error) {
	player, err := stores.PurchaseStore.GetPlayerRecord(*p.stores, t.ID)
	if (errors.Is(err, models.UserDoesNotExist{})) {
		return false, models.UserDoesNotExist{}
	}
	if l.ListerID == player.PlayerID {
		return false, models.UserBidOnOwnBid{}
	}
	if err != nil {
		return false, err
	}

	transaction := models.Transaction{
		ListingID: l.ListingID,
		Lister:    l.ListerID,
		Buyer:     t.ID,
		Price:     l.Price,
	}

	if l.Price > player.Clubcoin {
		return false, nil
	}
	err2 := stores.PurchaseStore.Purchase(*p.stores, transaction)
	if err2 != nil {
		return false, err2
	}
	return true, nil

}
