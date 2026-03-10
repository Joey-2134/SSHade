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

## Deploying to ECS (GitHub Actions)

Pushes to `main` trigger an automatic build and deploy: the workflow builds the Docker image, pushes it to Amazon ECR, and updates the ECS service.

### One-time setup

1. **ECR repository**  
   Create a repository (e.g. `sshade`) in ECR. Your ECS task definition should reference this image (e.g. `{account}.dkr.ecr.{region}.amazonaws.com/sshade:latest`).

2. **GitHub OIDC and IAM role**  
   - In **AWS IAM** → Identity providers, add an OIDC provider: `https://token.actions.githubusercontent.com` (no thumbprint, audience `sts.amazonaws.com`).  
   - Create an IAM role that trusts `token.actions.githubusercontent.com` with a condition on your repo (e.g. `repo:YourOrg/SSHade`).  
   - Attach policies (or inline policy) that allow:
     - **ECR**: `GetAuthorizationToken`; and for the repo resource, `BatchCheckLayerAvailability`, `GetDownloadUrlForLayer`, `BatchGetImage`, `PutImage`, `InitiateLayerUpload`, `UploadLayerPart`, `CompleteLayerUpload`.  
     - **ECS**: `ecs:UpdateService`, `ecs:DescribeServices`, and (if the role is used to register/update task definitions) `ecs:RegisterTaskDefinition`, `ecs:DescribeTaskDefinition`, `iam:PassRole` for the task execution role.  
   - Copy the role ARN (e.g. `arn:aws:iam::123456789012:role/github-actions-ecs`).

3. **GitHub repo configuration**  
   In the repo: **Settings → Secrets and variables → Actions**:
   - **Secrets**: add `AWS_ROLE_ARN` = the IAM role ARN from step 2.  
   - **Variables**: add  
     - `AWS_REGION` (e.g. `us-east-1`)  
     - `ECR_REPOSITORY` (e.g. `sshade`)  
     - `ECS_CLUSTER` (your ECS cluster name)  
     - `ECS_SERVICE` (your ECS service name)

After this, every push to `main` deploys the new image. You can also run the workflow manually from the **Actions** tab (“Deploy to ECS” → “Run workflow”).

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
