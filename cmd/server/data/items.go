package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreItem struct {
	GameID    int
	ItemID    string
	Cost      int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var storeItemsBaseSelect = "SELECT game_id, item_id, name, cost, created_at, updated_at FROM store_items "

func StoreItemsGet(ctx context.Context, pdb *pgxpool.Pool, storeItem *StoreItem) error {
	err := pdb.
		QueryRow(ctx, storeItemsBaseSelect+"WHERE game_id = $1 AND item_id = $2", storeItem.GameID, storeItem.ItemID).
		Scan(
			&storeItem.GameID,
			&storeItem.ItemID,
			&storeItem.Name,
			&storeItem.Cost,
			&storeItem.CreatedAt,
			&storeItem.UpdatedAt,
		)
	return err
}
