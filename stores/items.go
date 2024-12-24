package stores

import (
	"ClubmineStoreService/database"
	"ClubmineStoreService/models"
	"ClubmineStoreService/models/RequestTypes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	ItemStore interface {
		GetItemFromId(id string) (models.Item, error)
		GetItemFromStack(stack string) (models.Item, error)
		CreateItem(item *RequestTypes.CreateItemRequest) error
		DeleteItem(id string) error
	}

	itemStore struct{}
)

func (i itemStore) GetItemFromId(id string) (models.Item, error) {
	query :=
		`
		SELECT item_id, item_stack
		FROM items
		WHERE item_id = $1
		`
	var item models.Item
	if err := pgxscan.Get(context.Background(), database.DB, &item, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Item{}, models.NoRowsError{}
		}
		return models.Item{}, err
	}
	return item, nil
}

// GetItemFromStack does not work. character limit.
func (i itemStore) GetItemFromStack(stack string) (models.Item, error) {
	query :=
		`
		SELECT item_id, item_stack
		FROM items
		WHERE item_stack = $1
		`

	var item models.Item
	if err := pgxscan.Get(context.Background(), database.DB, &item, query, stack); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Item{}, models.NoRowsError{}
		}
		return models.Item{}, err
	}
	return item, nil
}

func (i itemStore) CreateItem(item *RequestTypes.CreateItemRequest) error {
	var count int
	if err := database.DB.QueryRow(context.Background(), countQuery, item.ItemID+"%").Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		var existingItemStack string
		err := database.DB.QueryRow(context.Background(), itemStackQuery, item.ItemStack).Scan(&existingItemStack)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		} else if existingItemStack == item.ItemStack {
			return models.AlreadyExistsInDatabaseError{}
		}
	}

	var maxVariant sql.NullInt64
	if err := database.DB.QueryRow(context.Background(), maxVariantQuery, item.ItemID+"%").Scan(&maxVariant); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	newVariant := 1
	if maxVariant.Valid {
		newVariant = int(maxVariant.Int64) + 1
	}

	newItemID := fmt.Sprintf("%s:%d", item.ItemID, newVariant)

	_, err := database.DB.Exec(context.Background(), insertQuery, newItemID, item.ItemStack, item.CreatorID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return models.ForeignKeyConstraint{}
			}
			if pgErr.Code == "23505" {
				return models.AlreadyExistsInDatabaseError{}
			}
		}
		return err
	}

	return nil
}

const (
	countQuery      = "SELECT COUNT(*) FROM items WHERE item_id LIKE $1"
	maxVariantQuery = "SELECT MAX(CAST(SPLIT_PART(item_id, ':', 3) AS INTEGER)) FROM items WHERE item_id LIKE $1"
	itemStackQuery  = "SELECT item_stack FROM items WHERE item_stack = $1"
	insertQuery     = `
		INSERT INTO items (item_id, item_stack, creator_id)
		VALUES($1, $2, $3)
	`
)

func (i itemStore) DeleteItem(id string) error {
	query := `
        DELETE FROM items
        WHERE item_id = $1
    `

	_, err := database.DB.Query(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
