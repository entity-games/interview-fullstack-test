package data

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Coins     int       `json:"coins"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LookupUser struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

var userBaseSelect = "SELECT user_id, username, coins, created_at, updated_at FROM users "
var userSearchSelect = "SELECT user_id, username FROM users "

func NewUserRegistration(ctx context.Context, pdb *pgxpool.Pool, user *User) (err error) {
	err = usersCreate(ctx, pdb, user)
	return
}

func usersCreate(ctx context.Context, pdb *pgxpool.Pool, user *User) error {
	user.UserID = uuid.New().String()
	user.Coins = 50
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt

	_, err := pdb.Exec(ctx, `
        INSERT INTO users (user_id, username, coins, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `, user.UserID, user.Username, user.Coins, user.CreatedAt, user.UpdatedAt)

	return err
}

func UsersGet(ctx context.Context, pdb *pgxpool.Pool, user *User) error {
	err := pdb.
		QueryRow(ctx, userBaseSelect+"WHERE user_id = $1", user.UserID).
		Scan(&user.UserID, &user.Username, &user.Coins, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func UsersGetByUsername(ctx context.Context, pdb *pgxpool.Pool, user *LookupUser) error {
	err := pdb.
		QueryRow(ctx, userSearchSelect+"WHERE username = $1", user.Username).
		Scan(&user.UserID, &user.Username)
	return err
}

func UsersUpdateCoins(ctx context.Context, pdb *pgxpool.Pool, user *User, delta int) error {
	user.Coins += delta
	user.UpdatedAt = time.Now().UTC()

	_, err := pdb.Exec(ctx, `
        UPDATE users SET coins = $1, updated_at = $2
        WHERE user_id = $3
    `, user.Coins, user.UpdatedAt, user.UserID)

	return err
}
