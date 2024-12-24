package stores

import (
	"ClubmineStoreService/database"
	"ClubmineStoreService/logger"
	"ClubmineStoreService/models"
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type (
	ListingStore interface {
		GetAll(pageNum int) ([]models.ListingResponse, error)
		GetAllByListerID(id string, pageNum int) ([]models.ListingResponse, error)
		GetListing(listingId string) (models.ListingResponse, error)
		CreateListing(listing *models.Listing) (bool, error)
		DeleteListing(listingId string) error
		GetItemFromID(itemId string) (models.Item, error)
		GetItemFromStack(baseId string, itemStack string) (models.Item, error)
	}

	listingStore struct{}
)

func (l listingStore) GetAll(pageNum int) ([]models.ListingResponse, error) {
	query := `
        SELECT * FROM listings WHERE listing_id 
        IN (
            SELECT listing_id
            FROM (
                SELECT listing_id, ROW_NUMBER() OVER (ORDER BY price DESC) AS row_num
                FROM listings
            ) AS ranked_listings
            WHERE row_num <= $1 + 1 * 28
            AND row_num > $1 * 28
        )
        ORDER BY price DESC;
    `

	rows, err := database.DB.Query(context.Background(), query, pageNum)
	if err != nil {
		return []models.ListingResponse{}, err
	}
	defer rows.Close()

	var listings []models.ListingResponse
	err = pgxscan.ScanAll(&listings, rows)
	if err != nil || len(listings) == 0 {
		return []models.ListingResponse{}, err
	}
	for i := range listings {
		item, err := l.GetItemFromID(listings[i].ItemID)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		listings[i].ItemStack = item.ItemStack
	}

	return listings, nil
}

func (l listingStore) GetAllByListerID(id string, pageNum int) ([]models.ListingResponse, error) {
	query := `
        SELECT * FROM listings WHERE lister_id = $1
        ORDER BY price DESC
        LIMIT 28 OFFSET $2 * 28;
    `

	rows, err := database.DB.Query(context.Background(), query, id, pageNum)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []models.ListingResponse
	err = pgxscan.ScanAll(&listings, rows)
	if err != nil {
		return nil, err
	}
	for i := range listings {
		item, err := l.GetItemFromID(listings[i].ItemID)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		listings[i].ItemStack = item.ItemStack
	}

	return listings, nil
}

func (l listingStore) GetListing(listingId string) (models.ListingResponse, error) {
	query := `
        SELECT listing_id, lister_id, item_id, price, created_at
        FROM listings
        WHERE listing_id = $1
    `

	var listing models.ListingResponse
	if err := pgxscan.Get(context.Background(), database.DB, &listing, query, listingId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ListingResponse{}, nil
		}
		return models.ListingResponse{}, err
	}
	item, err := l.GetItemFromID(listing.ItemID)
	if err != nil {
		logger.Fatal(err.Error())
	}
	listing.ItemStack = item.ItemStack
	return listing, nil
}

func (l listingStore) CreateListing(listing *models.Listing) (bool, error) {
	query := `
        INSERT INTO listings (listing_id, lister_id, item_id, price, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `
	if database.DB == nil {
		logger.Fatal("sum went wrong")
	}

	variantObject, err2 := l.GetItemFromStack(listing.ItemID, listing.ItemStack)
	if err2 != nil {
		if errors.Is(err2, models.ItemDoesNotExist{}) {
			return false, err2
		} else {
			return false, err2
		}
	}

	listing.ItemID = variantObject.ItemID

	_, err := database.DB.Exec(context.Background(), query, listing.ListingID, listing.ListerID, listing.ItemID, listing.Price, listing.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return false, models.ForeignKeyConstraint{}
			}
		}
		logger.Error("Something went wrong : ", zap.Error(err))
		return false, err
	}
	return true, nil
}

func (l listingStore) DeleteListing(listingId string) error {
	query := `
        DELETE FROM listings
        WHERE listing_id = $1
    `

	_, err := database.DB.Query(context.Background(), query, listingId)
	if err != nil {
		return err
	}
	return nil
}

func (l listingStore) GetItemFromID(itemId string) (models.Item, error) {
	query :=
		`
		SELECT item_id, item_stack
		FROM items
		WHERE item_id = $1
		`

	var item models.Item
	if err := pgxscan.Get(context.Background(), database.DB, &item, query, itemId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Item{}, nil
		}
		return models.Item{}, err
	}
	return item, nil
}

func (l listingStore) GetItemFromStack(baseId string, itemStack string) (models.Item, error) {
	query :=
		`
		SELECT item_id, item_stack
		FROM items
		WHERE item_stack = $1 AND item_id LIKE $2 || '%'
		`

	var item models.Item
	if err := pgxscan.Get(context.Background(), database.DB, &item, query, itemStack, baseId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Item{}, models.ItemDoesNotExist{}
		}
		return models.Item{}, err
	}
	return item, nil
}
