package db

import (
	"database/sql"
)

type User struct {
	ID           int
	Username     string
	Fingerprint  string
	FactionID    sql.NullInt64
	LastPlacedAt sql.NullTime
	CreatedAt    sql.NullTime
}

func GetUserByFingerprint(db *sql.DB, fingerprint string) (*User, error) {
	row := db.QueryRow("SELECT id, username, ssh_key_fingerprint, faction_id, last_placed_at, created_at FROM users WHERE ssh_key_fingerprint = ?", fingerprint)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Fingerprint, &user.FactionID, &user.LastPlacedAt, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
