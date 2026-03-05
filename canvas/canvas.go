package canvas

import (
	"context"
	"database/sql"
	"sync"

	"github.com/Joey-2134/SSHade/db"
)

const defaultColour = "#cccccc"

// Canvas is the shared in-memory grid of pixels, persisted to SQLite.
type Canvas struct {
	mu     sync.RWMutex
	width  int
	height int
	Pixels [][]Pixel
}

// New allocates a canvas of the given dimensions and fills it with the default colour.
func New(width, height int) *Canvas {
	pixels := make([][]Pixel, height)
	for y := range height {
		pixels[y] = make([]Pixel, width)
		for x := range width {
			pixels[y][x] = Pixel{X: x, Y: y, ColourHex: defaultColour}
		}
	}
	return &Canvas{
		width:  width,
		height: height,
		Pixels: pixels,
	}
}

// LoadFromDB fills the canvas from the database. Pixels in the DB are applied;
// any cell not in the DB keeps the default colour from New.
func (c *Canvas) LoadFromDB(ctx context.Context, database *sql.DB) error {
	rows, err := db.LoadPixels(ctx, database)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range rows {
		if p.X >= 0 && p.X < c.width && p.Y >= 0 && p.Y < c.height {
			c.Pixels[p.Y][p.X] = Pixel{X: p.X, Y: p.Y, ColourHex: p.ColourHex}
		}
	}
	return nil
}

// PixelAt returns the pixel at (x, y). The second return is false if out of bounds.
func (c *Canvas) PixelAt(x, y int) (Pixel, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return Pixel{}, false
	}
	return c.Pixels[y][x], true
}

// Set updates the pixel at (x, y) in memory and persists it via the database.
// Returns an error if out of bounds or if the DB write fails.
func (c *Canvas) Set(ctx context.Context, database *sql.DB, x, y int, colourHex string) error {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return nil // skip out-of-bounds
	}
	c.mu.Lock()
	c.Pixels[y][x] = Pixel{X: x, Y: y, ColourHex: colourHex}
	c.mu.Unlock()
	return db.SetPixel(ctx, database, x, y, colourHex, nil, nil)
}

// Width returns the canvas width.
func (c *Canvas) Width() int { return c.width }

// Height returns the canvas height.
func (c *Canvas) Height() int { return c.height }
