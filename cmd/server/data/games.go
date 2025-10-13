package data

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GameUserContext struct {
	Unlocked bool `json:"unlocked"`
}

type Game struct {
	GameID      int             `json:"-"`
	Title       string          `json:"title"`
	TitleStub   string          `json:"-"`
	Description string          `json:"-"`
	Cost        int             `json:"-"`
	Category    string          `json:"-"`
	HomeUrl     string          `json:"-"`
	ApiKey      string          `json:"-"`
	CreatedAt   time.Time       `json:"-"`
	UpdatedAt   time.Time       `json:"-"`
	UserContext GameUserContext `json:"-"`
}

var gamesBaseSelect = "SELECT game_id, title, title_stub, description, cost, category, home_url, api_key, created_at, updated_at FROM games "

func GamesList(ctx context.Context, pdb *pgxpool.Pool, userUnlockedGameIds []int) <-chan Game {
	out := make(chan Game)

	go func() {
		defer close(out)

		conn, err := pdb.Acquire(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Connection error: ", "err", err)
			return
		}
		defer conn.Release()

		rows, err := conn.Query(ctx, gamesBaseSelect)
		if err != nil {
			slog.ErrorContext(ctx, "Query error: ", "err", err)
			return
		}
		defer rows.Close()

		var game Game
		game.UserContext = GameUserContext{}

		for rows.Next() {
			err := rows.Scan(
				&game.GameID,
				&game.Title,
				&game.TitleStub,
				&game.Description,
				&game.Cost,
				&game.Category,
				&game.HomeUrl,
				&game.ApiKey,
				&game.CreatedAt,
				&game.UpdatedAt,
			)
			if err != nil {
				slog.ErrorContext(ctx, "Scan error: ", "err", err)
				return
			}
			game.UserContext.Unlocked = slices.Contains(userUnlockedGameIds, game.GameID)
			out <- game
		}
	}()

	return out
}

func GamesGet(ctx context.Context, pdb *pgxpool.Pool, game *Game) error {
	err := pdb.
		QueryRow(ctx, gamesBaseSelect+"WHERE game_id = $1", game.GameID).
		Scan(&game.GameID, &game.Title, &game.TitleStub, &game.Description, &game.Cost, &game.Category, &game.HomeUrl, &game.ApiKey, &game.CreatedAt, &game.UpdatedAt)
	return err
}
