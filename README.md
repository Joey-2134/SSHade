# SSHade

A multiplayer pixel canvas accessible entirely over SSH. Connect, join a faction, and place pixels on a shared canvas that updates in real time for every connected user. No browser. No app. Just a terminal.

Inspired by [r/place](https://reddit.com/r/place) and [Terminal.shop](https://terminal.shop).

```
ssh sshade.example.com
```

## What it is

SSHade is a persistent, collaborative pixel art canvas rendered entirely in your terminal over SSH. Every connected user sees the same canvas. Place one pixel at a time, subject to a cooldown. Coordinate with your faction to claim territory before rival factions do.

## Features

- **Real-time multiplayer** - pixels placed by other users appear on your canvas instantly
- **Factions** - join a team, pick a colour, fight for territory
- **Cooldowns** - one pixel at a time, enforced server-side
- **Persistent canvas** - state survives server restarts via SQLite
- **Pure SSH** - no web app, no client to install; `ssh` is the only dependency

## Tech Stack

| Component | Library |
|-----------|---------|
| Language | Go |
| SSH server | [wish](https://github.com/charmbracelet/wish) (Charm) |
| Terminal UI | [bubbletea](https://github.com/charmbracelet/bubbletea) |
| Styling | [lipgloss](https://github.com/charmbracelet/lipgloss) |
| Database | SQLite via `modernc.org/sqlite` (pure Go, no CGO) |

## Project Structure

```
SSHade/
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
│   └── config.go         # Canvas size, cooldown duration, port, etc.
└── go.mod
```

## Running Locally

```bash
git clone https://github.com/your-username/SSHade
cd SSHade
go run .
```

Then in another terminal:

```bash
ssh localhost -p 2222
```

## Configuration

All values are configurable in `config/config.go`:

- Canvas width and height
- Cooldown duration per pixel placement
- SSH listen port
- SQLite file path
- Max connections per IP
- Faction definitions (name + colour)

## Requirements

- Go 1.21+
- A terminal with 256-colour support (24-bit recommended)

## Roadmap

- **Phase 1** - Core: SSH server, shared canvas, real-time updates
- **Phase 2** - Identity: usernames, SSH key auth, factions, cooldowns, leaderboard
- **Phase 3** - Maps: shaped canvases, scheduled resets, replay mode
- **Phase 4** - Stretch: cooldown economy, server events, web viewer
