package db

import (
	"context"
	"database/sql"
)

type Pixel struct {
	X         int
	Y         int
	ColourHex string
}

// LoadPixels returns all pixels from canvas_pixels, for populating the in-memory canvas on startup.
func LoadPixels(ctx context.Context, db *sql.DB) ([]Pixel, error) {
	rows, err := db.QueryContext(ctx, `SELECT x, y, colour_hex FROM canvas_pixels`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Pixel
	for rows.Next() {
		var p Pixel
		if err := rows.Scan(&p.X, &p.Y, &p.ColourHex); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// SetPixel writes a pixel to the canvas_pixels table (replace if exists) and appends to pixel_history.
func SetPixel(ctx context.Context, db *sql.DB, x, y int, colourHex string, factionID, placedBy *int64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx,
		`INSERT INTO canvas_pixels (x, y, colour_hex, faction_id, placed_by, placed_at)
		 VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(x, y) DO UPDATE SET colour_hex = ?, faction_id = ?, placed_by = ?, placed_at = CURRENT_TIMESTAMP`,
		x, y, colourHex, factionID, placedBy, colourHex, factionID, placedBy)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO pixel_history (x, y, colour_hex, faction_id, placed_by, placed_at)
		 VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		x, y, colourHex, factionID, placedBy)
	if err != nil {
		return err
	}
	return tx.Commit()
}
