package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const OrderStatusCreated int = 1
const OrderStatusPaid int = 2
const OrderStatusApplied int = 3

const OrderTypeAccess string = "access"
const OrderTypeGameGoods string = "ingame"
const OrderTypePlatformGoods string = "platform"

type Order struct {
	OrderID   string
	UserID    string
	GameID    int
	ItemID    string
	Status    int
	Type      string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ordersBaseSelect = "SELECT order_id, user_id, item_id, game_id, status, type, content, created_at, updated_at FROM orders "

func OrdersCreate(ctx context.Context, pdb *pgxpool.Pool, order *Order) error {
	order.OrderID = uuid.New().String()
	order.CreatedAt = time.Now().UTC()
	order.UpdatedAt = order.CreatedAt

	_, err := pdb.Exec(ctx, `
        INSERT INTO orders (order_id, user_id, item_id, game_id, status, type, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `,
		order.OrderID,
		order.UserID,
		order.ItemID,
		order.GameID,
		order.Status,
		order.Type,
		order.Content,
		order.CreatedAt,
		order.UpdatedAt,
	)

	return err
}

func OrdersUpdate(ctx context.Context, pdb *pgxpool.Pool, order *Order) error {
	order.UpdatedAt = time.Now().UTC()

	_, err := pdb.Exec(ctx, "UPDATE orders SET status = $1, updated_at = $2 WHERE order_id = $3",
		order.Status,
		order.UpdatedAt,
		order.OrderID,
	)

	return err
}

func OrdersGet(ctx context.Context, pdb *pgxpool.Pool, order *Order) error {
	err := pdb.
		QueryRow(ctx, ordersBaseSelect+"WHERE order_id = $1", order.OrderID).
		Scan(
			&order.OrderID,
			&order.UserID,
			&order.ItemID,
			&order.GameID,
			&order.Status,
			&order.Type,
			&order.Content,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
	return err
}

func ApplyOrder(ctx context.Context, pdb *pgxpool.Pool, order *Order) error {
	if order.Type == OrderTypeAccess {
		// unlock the game
		gameId, _ := strconv.Atoi(order.ItemID)
		gameData := &GameData{UserID: order.UserID, GameID: gameId}
		err := GameDataGet(ctx, pdb, gameData)
		if errors.Is(err, sql.ErrNoRows) {
			err = GameDataUpdate(ctx, pdb, gameData)
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("Unable to update gamedata: %s", err.Error()))
				return err
			}
		}

		// notify the game about the purchase of the game

	} else if order.Type == OrderTypePlatformGoods {
		// modify platform data
		user := &User{UserID: order.UserID}
		err := UsersGet(ctx, pdb, user)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Unable find user: %s", err.Error()))
			return err
		}

		// TODO: start utilizing StoreItem.Content and Order.Content
		parts := strings.Split(order.ItemID, ":")
		delta, _ := strconv.Atoi(parts[1])

		err = UsersUpdateCoins(ctx, pdb, user, delta)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Unable to update user data: %s", err.Error()))
			return err
		}
	} else if order.Type == OrderTypeGameGoods {
		// notify the game about the purchased goods

	} else {
		// unsupported
		return errors.New("invalid order type")
	}

	return nil
}
