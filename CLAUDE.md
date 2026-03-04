# SSH Pixel Canvas - Project Brief

## Concept

A multiplayer pixel canvas accessible entirely over SSH. Users connect, see the current state of a shared canvas rendered in their terminal, choose a faction, and place one pixel at a time subject to a cooldown. The canvas is persistent and shared across all connected users in real time.

Inspired by Reddit's r/place and Primagen's Terminal.shop. The goal is a polished, absurd, fully committed SSH-native experience.

## Tech Stack

- **Language**: Go
- **SSH server**: [wish](https://github.com/charmbracelet/wish) by Charm
- **Terminal UI**: [bubbletea](https://github.com/charmbracelet/bubbletea) for the TUI framework
- **Styling**: [lipgloss](https://github.com/charmbracelet/lipgloss) for terminal colour and layout
- **Database**: SQLite via `modernc.org/sqlite` (pure Go, no CGO required)
- **Canvas broadcast**: In-process pub/sub or Go channels for real time updates across sessions

## Project Phases

### Phase 1 - Core (Start Here)
Get a working SSH server where multiple users can connect simultaneously and place pixels on a shared canvas that updates in real time for all connected users.

Milestone: Two terminal windows SSH'd into localhost, placing a pixel in one immediately renders in the other.

### Phase 2 - Identity and Factions
- Username assignment on connect (prompt on first connect, remembered by SSH key)
- Faction selection UI shown to new users displaying current canvas state to inform their choice
- Faction colours and territory tracking
- Cooldown enforcement per user stored in SQLite
- Leaderboard showing faction territory percentages

### Phase 3 - Maps and Resets
- Custom shaped canvases (world map, country outlines) instead of open grid
- Map rotation on a configurable schedule (daily/weekly reset)
- Historical snapshots saved on each reset
- Canvas replay mode - scrub through the edit history of a session

### Phase 4 - Stretch
- Cooldown economy (earn credits over time, spend to place rectangles or protect regions)
- Procedural server events (voids, challenges, bonus windows)
- Web viewer for the canvas (read-only, for people without a terminal)

## Architecture Notes

### SSH Session Lifecycle
Each SSH connection gets its own bubbletea program instance. Session state (username, faction, last placement time) is loaded from SQLite on connect. The bubbletea model handles input and renders the canvas for that session.

### Canvas State
Canvas is a 2D array of pixels stored in memory as the source of truth, persisted to SQLite. On server start, load canvas from SQLite into memory. All reads serve from memory. All writes go to memory first then async to SQLite.

### Real Time Updates
Use a central broadcaster - a Go channel or simple pub/sub - that all active sessions subscribe to. When any user places a pixel, the change is pushed to the broadcaster which fans it out to every connected session, triggering a re-render.

### Concurrency
Canvas writes must be mutex-protected. Two simultaneous placements at the same coordinate resolve as last-write-wins. Cooldown validation must be server-side only - never trust client state for this.

## Key Implementation Concerns

**Terminal size variance** - Query terminal dimensions on connect using wish/bubbletea hooks. Handle resize events. Define a minimum canvas display size and show a clear error if the terminal is too small rather than rendering garbage.

**Colour depth variance** - Not all terminals support 24-bit colour. Decide on a minimum (256 colour is a reasonable baseline) and document it. Lipgloss handles colour fallback to some degree but test on multiple terminals.

**Cooldown enforcement** - Store `last_placed_at` timestamp per user in SQLite. Validate server-side on every placement attempt. Return remaining cooldown time to the client for display.

**Auth** - Accept both password and SSH key auth initially. Identify users by SSH key fingerprint where available, fall back to username/password. Key fingerprint is more stable as a user identity anchor.

**Bot mitigation** - Rate limit connections per IP at the wish middleware layer. This does not need to be complex initially, just a connection rate limit.

**Concurrent writes** - Use a `sync.RWMutex` on the in-memory canvas. Read lock for renders, write lock for placements.

## Database Schema (Initial)

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    ssh_key_fingerprint TEXT,
    faction_id INTEGER,
    last_placed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE factions (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    colour_hex TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE canvas_pixels (
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    colour_hex TEXT NOT NULL,
    faction_id INTEGER,
    placed_by INTEGER,
    placed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (x, y)
);

CREATE TABLE pixel_history (
    id INTEGER PRIMARY KEY,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    colour_hex TEXT NOT NULL,
    faction_id INTEGER,
    placed_by INTEGER,
    placed_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Project Structure

```
ssh-canvas/
├── main.go               # Entry point, SSH server setup
├── canvas/
│   ├── canvas.go         # In-memory canvas state, mutex, broadcaster
│   └── pixel.go          # Pixel type definitions
├── db/
│   ├── db.go             # SQLite connection and migrations
│   ├── users.go          # User queries
│   ├── pixels.go         # Canvas persistence queries
│   └── factions.go       # Faction queries
├── session/
│   ├── session.go        # Per-connection state
│   └── handler.go        # Wish middleware and session init
├── ui/
│   ├── model.go          # Bubbletea model (main app state)
│   ├── canvas_view.go    # Canvas rendering logic
│   ├── faction_view.go   # Faction selection screen
│   └── hud.go            # Cooldown timer, faction info overlay
├── config/
│   └── config.go         # Canvas size, cooldown duration, port etc
└── go.mod
```

## Configuration (config.go)

Make these values configurable from the start rather than hardcoded:
- Canvas width and height
- Cooldown duration per placement
- SSH listen port
- SQLite file path
- Max connections per IP
- Faction definitions (name + colour)

## Development Order

1. `main.go` + bare wish SSH server that accepts connections and prints hello
2. Static canvas render in bubbletea - hardcoded colours, no interaction
3. Cursor movement across the canvas
4. SQLite setup and canvas persistence (load on start, save on write)
5. Pixel placement for a single user
6. Broadcaster - second connected session sees placements from the first in real time
7. User identity - prompt for username on first connect, store against SSH key
8. Cooldown - enforce and display remaining time
9. Faction selection screen
10. Faction territory tracking and leaderboard

Do not move to phase 2 features until step 6 is working solidly.

## References

- [Wish documentation](https://github.com/charmbracelet/wish)
- [Bubbletea documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss documentation](https://github.com/charmbracelet/lipgloss)
- [Terminal.shop](https://terminal.shop) - reference for SSH app UX
- [r/place](https://reddit.com/r/place) - original concept reference
