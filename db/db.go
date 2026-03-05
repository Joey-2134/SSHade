package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const DefaultPath = "./sshade.db"

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite %q: %w", path, err)
	}
	if err := Migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	stmt := `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    ssh_key_fingerprint TEXT,
    faction_id INTEGER,
    last_placed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS factions (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    colour_hex TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS canvas_pixels (
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    colour_hex TEXT NOT NULL,
    faction_id INTEGER,
    placed_by INTEGER,
    placed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (x, y)
);

CREATE TABLE IF NOT EXISTS pixel_history (
    id INTEGER PRIMARY KEY,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    colour_hex TEXT NOT NULL,
    faction_id INTEGER,
    placed_by INTEGER,
    placed_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`
	_, err := db.Exec(stmt)
	return err
}
