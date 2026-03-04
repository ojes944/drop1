package db

import (
	"context"
	"time"
)

func CreateUser(email, password, name string) error {
	_, err := DB.ExecContext(context.Background(), "INSERT INTO users (email, password, name) VALUES ($1, $2, $3)", email, password, name)
	return err
}

func GetUserByEmail(email string) (id int, password string, name string, err error) {
	row := DB.QueryRowContext(context.Background(), "SELECT id, password, name FROM users WHERE email=$1", email)
	err = row.Scan(&id, &password, &name)
	return
}

func StoreResetToken(userID int, token string, expiresAt time.Time) error {
	_, err := DB.ExecContext(context.Background(), "INSERT INTO password_resets (user_id, token, expires_at) VALUES ($1, $2, $3)", userID, token, expiresAt)
	return err
}

func GetResetTokenDB(userID int) (token string, expiresAt time.Time, err error) {
	row := DB.QueryRowContext(context.Background(), "SELECT token, expires_at FROM password_resets WHERE user_id=$1 ORDER BY expires_at DESC LIMIT 1", userID)
	err = row.Scan(&token, &expiresAt)
	return
}
