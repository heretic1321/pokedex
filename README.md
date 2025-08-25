## Pokedex CLI

Smoll, command‑line Pokedex written in Go. Browse regions, explore areas, catch Pokémon, and inspect what you’ve caught. It talks to the public PokeAPI.

Planned fun extras: ASCII art sprites, simple battles, and more tinkering.

### Features
- **map/mapb**: page through location areas (20 at a time)
- **explore <area>**: list Pokémon that can appear in an area (e.g., `explore canalave-city`)
- **catch <name>**: try to catch a Pokémon by name
- **inspect <name>**: view height, weight, and base stats for a caught Pokémon
- **pokedex**: list all Pokémon you’ve caught
- **help/exit**: show help or quit

### Requirements
- **Go**: latest stable Go (recommended). Any recent Go toolchain should work.
- **Internet access**: uses `https://pokeapi.co/api/v2/` (no API key required).
- Optional: **make** for convenience targets (`make build`, `make run`).

### Install
Using Go (installs to your `$GOBIN` or `$(go env GOPATH)/bin`):

```bash
go install github.com/heretic1321/pokedex/cmd/pokedex@latest
```

Then run:

```bash
pokedex
```

Path tips:
- Linux/macOS: ensure `$(go env GOPATH)/bin` is in `PATH`.
- Windows (PowerShell): add `%USERPROFILE%\go\bin` to your PATH.
- Windows (Git Bash): `export PATH="$PATH:$(go env GOPATH)/bin"`.

### Build from source
Clone and build:

```bash
git clone https://github.com/heretic1321/pokedex.git
cd pokedex
make build        # or: go build -o pokedex ./cmd/pokedex
./pokedex         # on Windows: ./pokedex.exe
```

Quick run without building a binary:

```bash
make run          # or: go run ./cmd/pokedex
```

### Usage
Start the CLI and type commands at the `Pokedex >` prompt. Examples:

```bash
pokedex

# List areas (next 20)
Pokedex > map

# Go back to previous page of areas
Pokedex > mapb

# Explore an area to see encounterable Pokémon
Pokedex > explore canalave-city

# Try to catch a Pokémon by name
Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!

# Show what you’ve caught
Pokedex > pokedex

# Inspect stats of a caught Pokémon
Pokedex > inspect pikachu

# See available commands
Pokedex > help

# Quit
Pokedex > exit
```

Notes:
- Area and Pokémon names are lowercase and hyphenated as per PokeAPI (e.g., `canalave-city`). Use `map` to discover area names first.
- Catch chance is based on the Pokémon’s base experience; sometimes they escape—try again!

### Behavior & internals (practical details)
- **API**: `https://pokeapi.co/api/v2/`
- **Caching**: in‑memory cache with ~10s TTL to reduce repeat network calls.
- **Pagination**: `map`/`mapb` move forward/back by 20 areas per page.

### Development
Run locally:

```bash
go run ./cmd/pokedex
```

Run tests (includes cache behavior tests):

```bash
go test ./...
```

Project layout highlights:
- `cmd/pokedex`: entrypoint (`main.go`)
- `internal/cli`: interactive shell and command handlers
- `internal/pokedex`: service layer, in‑memory state and cache usage
- `pkg/pokeapi`: thin client and response models for PokeAPI

### Troubleshooting
- "command not found": ensure Go’s bin directory is on your PATH (see Install section).
- Network errors or slow responses: PokeAPI may be rate‑limited or down; try again.
- "pokemon name invalid" or "no area specified": check spelling; use `map`/`help`.

### Roadmap (ideas)
- ASCII art for Pokémon
- Lightweight battle/battle‑sim commands
- More filters and search commands

Enjoy!


