package stores

import (
	"ClubmineStoreService/database"
	"ClubmineStoreService/logger"
	"ClubmineStoreService/models"
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type (
	PurchaseStore interface {
		Purchase(transaction models.Transaction) error
		GetPlayerRecord(id string) (models.Player, error)
	}

	purchaseStore struct{}
)

func (p purchaseStore) Purchase(transaction models.Transaction) error {

	// $1 = PRICE
	// $2 = BUYER_ID
	// $3 = LISTER_ID
	// $4 = LISTING_ID

	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			logger.Error("Rolling back transaction due to error", zap.Error(err))
			err := tx.Rollback(context.Background())
			if err != nil {
				return
			}
			return
		}

		logger.Info("Committing transaction")
		if err := tx.Commit(context.Background()); err != nil {
			logger.Error("Failed to commit transaction", zap.Error(err))
		}
	}()

	_, err = tx.Exec(context.Background(), `DELETE FROM listings WHERE listing_id = $1`, transaction.ListingID)
	if err != nil {
		logger.Error("Error deleting listing", zap.Error(err))
		return err
	}

	_, err = tx.Exec(context.Background(), `UPDATE players SET clubcoin = clubcoin - $1 WHERE player_id = $2`, transaction.Price, transaction.Buyer)
	if err != nil {
		logger.Error("Error updating buyer's balance", zap.Error(err))
		return err
	}

	_, err = tx.Exec(context.Background(), `UPDATE players SET clubcoin = clubcoin + $1 WHERE player_id = $2`, transaction.Price, transaction.Lister)
	if err != nil {
		logger.Error("Error updating seller's balance", zap.Error(err))
		return err
	}

	_, err = tx.Exec(context.Background(), `INSERT INTO transactions (listing_id, lister_id, buyer_id, price) VALUES ($1, $2, $3, $4)`, transaction.ListingID, transaction.Lister, transaction.Buyer, transaction.Price)
	if err != nil {
		logger.Error("Error inserting transaction record", zap.Error(err))
		return err
	}

	return nil
}

func (p purchaseStore) GetPlayerRecord(id string) (models.Player, error) {
	query := `SELECT player_id, clubcoin FROM players WHERE player_id = $1`

	var player models.Player
	if err := pgxscan.Get(context.Background(), database.DB, &player, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Player{}, models.UserDoesNotExist{}
		}
		return models.Player{}, err
	}
	return player, nil
}
