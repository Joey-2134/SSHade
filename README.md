# SSHade

A multiplayer pixel canvas accessible entirely over SSH. Connect, join a faction, and place pixels on a shared canvas that updates in real time for every connected user. No browser. No app. Just a terminal.

This was a personal project made with the aim of learning to make ssh programs and styled terminal UI's and at the same time experiment with Cursor AI.

Inspired by [r/place](https://reddit.com/r/place) and [Terminal.shop](https://terminal.shop).

```
ssh sshade.net
```

## What it is

SSHade is a persistent, collaborative pixel art canvas rendered entirely in your terminal over SSH. Every connected user sees the same canvas. Place one pixel at a time.

## Features

- **Real-time multiplayer** - pixels placed by other users appear on your canvas instantly
- **Factions** - join a team, pick a colour
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
├── main.go                    # Entry point, SSH server setup
├── go.mod
├── go.sum
├── Dockerfile
├── banner.txt                 # ASCII banner shown on connect
├── .gitignore
├── .github/
│   └── workflows/
│       └── deploy.yml         # CI/CD
├── canvas/
│   ├── canvas.go              # In-memory canvas state, mutex
│   ├── pixel.go               # Pixel type definitions
│   └── broadcast.go           # Real-time updates to connected sessions
├── constants/
│   ├── keymap.go              # Key bindings
│   └── uiconstants.go         # UI layout constants
├── db/
│   ├── db.go                  # SQLite connection and migrations
│   ├── users.go               # User queries
│   ├── pixels.go              # Canvas persistence queries
│   ├── factions.go            # Faction queries
│   └── migrations/
├── ui/
│   ├── components/            # Reusable UI pieces
│   │   ├── header.go
│   │   ├── footer.go
│   │   ├── grid.go            # Canvas grid rendering
│   │   ├── emptyfactions.go
│   │   └── factioncreation.go
│   └── screens/               # Per-screen Bubbletea models
│       ├── model.go           # Main app state, screen routing
│       ├── splash.go          # Welcome / loading
│       ├── usercreation.go    # Username prompt
│       ├── factionselection.go
│       ├── factioncreation.go
│       └── keys_ssh.go        # SSH connection hint
└── CLAUDE.md                  # Project brief for AI tooling
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

## Requirements

- Go 1.21+
- A terminal with 256-colour support (24-bit recommended)