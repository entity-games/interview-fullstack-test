package data

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameData struct {
	UserID    string
	GameID    int
	Data      map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var gamedataBaseSelect = "SELECT user_id, game_id, data, created_at, updated_at FROM game_data "

func GameDataUpdate(ctx context.Context, pdb *pgxpool.Pool, gamedata *GameData) error {
	gamedata.UpdatedAt = time.Now().UTC()
	if gamedata.CreatedAt.IsZero() {
		gamedata.CreatedAt = gamedata.UpdatedAt
	}

	jsonData, err := json.Marshal(&gamedata.Data)
	if err != nil {
		return err
	}

	_, err = pdb.Exec(ctx, `
		INSERT INTO game_data (user_id, game_id, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) 
		ON CONFLICT (user_id, game_id) 
		    DO UPDATE SET data = EXCLUDED.data, updated_at = EXCLUDED.updated_at`,
		gamedata.UserID, gamedata.GameID, jsonData, gamedata.CreatedAt, gamedata.UpdatedAt)

	return err
}

func GameDataGet(ctx context.Context, pdb *pgxpool.Pool, gamedata *GameData) error {
	gamedata.Data = make(map[string]string)
	var jsonData []byte

	err := pdb.
		QueryRow(ctx, gamedataBaseSelect+"WHERE user_id = $1 and game_id = $2", gamedata.UserID, gamedata.GameID).
		Scan(&gamedata.UserID, &gamedata.GameID, &jsonData, &gamedata.CreatedAt, &gamedata.UpdatedAt)

	if jsonData != nil && len(jsonData) > 2 {
		err = json.Unmarshal(jsonData, &gamedata.Data)
		if err != nil {
			return err
		}
	}

	return err
}

func GameDataGetUnlocked(ctx context.Context, pdb *pgxpool.Pool, userId string) []int {
	var unlockedGames []int

	rows, err := pdb.Query(ctx, "SELECT game_id FROM game_data WHERE user_id=$1", userId)
	if err != nil {
		slog.ErrorContext(ctx, "Query error: ", "err", err)
		return unlockedGames
	}
	defer rows.Close()

	var gameId int
	_, err = pgx.ForEachRow(rows, []any{&gameId}, func() error {
		unlockedGames = append(unlockedGames, gameId)
		return nil
	})

	return unlockedGames
}
