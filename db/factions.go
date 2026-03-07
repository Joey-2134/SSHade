package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Faction struct {
	ID        int
	Name      string
	ColourHex string
	CreatedAt sql.NullTime
}

func GetAllFactions(ctx context.Context, db *sql.DB) ([]Faction, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, name, colour_hex, created_at FROM factions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var factions []Faction
	for rows.Next() {
		var faction Faction
		if err := rows.Scan(&faction.ID, &faction.Name, &faction.ColourHex, &faction.CreatedAt); err != nil {
			return nil, err
		}
		factions = append(factions, faction)
	}
	return factions, rows.Err()
}

func GetFactionByID(ctx context.Context, db *sql.DB, id int) (Faction, error) {
	row := db.QueryRowContext(ctx, "SELECT id, name, colour_hex, created_at FROM factions WHERE id = ?", id)
	var faction Faction
	if err := row.Scan(&faction.ID, &faction.Name, &faction.ColourHex, &faction.CreatedAt); err != nil {
		return Faction{}, err
	}
	return faction, nil
}

func CreateFaction(ctx context.Context, db *sql.DB, name string, colourHex string) (Faction, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return Faction{}, err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, "INSERT INTO factions (name, colour_hex) VALUES (?, ?)", name, colourHex)
	if err != nil {
		return Faction{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Faction{}, err
	}
	if err := tx.Commit(); err != nil {
		return Faction{}, err
	}
	return GetFactionByID(ctx, db, int(id))
}

func (f Faction) String() string {
	return fmt.Sprintf("%s (%s)", f.Name, f.ColourHex)
}

func (f Faction) GetFactionName() string {
	return f.Name
}
