# Dwaar

Dwaar is a powerful and intuitive Terminal User Interface (TUI) for managing SSH and local port forwarding rules. Built in Go using the renowned [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework, it brings an Elm-like architecture to terminal applications, ensuring a robust and responsive user experience.

## Overview

Managing port forwarding rules via CLI can become cumbersome, especially with multiple connections and SSH hops. Dwaar simplifies this by providing a unified interface to:
- Configure local port to remote host:port mappings.
- Start and stop SSH tunnels seamlessly.
- View real-time connection logs and tunnel statuses.
- Persist configurations for quick setups on subsequent runs.

## Features

- **TUI Management**: Add, edit, and delete port forwarding rules using a built-in form interface.
- **Tunnel Control**: Start, stop, and restart individual tunnels with single keystrokes.
- **Status Indicators**: Visual cues for `running`, `stopped`, `connecting`, and `error` states.
- **Persistent Configuration**: Rules are saved to and loaded from a YAML configuration file.
- **SSH Support**: Forward ports over SSH with key and password authentication.

## Tech Stack

- **Core Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles)
- **Forms**: [Huh](https://github.com/charmbracelet/huh)
- **SSH Engine**: `golang.org/x/crypto/ssh`

## Installation

Ensure you have [Go](https://golang.org/doc/install) installed (version 1.21+ recommended).

```bash
git clone <repository-url>/dwaar
cd dwaar
go build -o dwaar ./cmd/dwaar
```

## Usage

Run the built executable to start the application:

```bash
./dwaar
```

### Key Bindings

| Key | Action |
|---|---|
| `↑` / `↓` or `j` / `k` | Navigate list |
| `Enter` | View details / open logs |
| `n` | New rule |
| `e` | Edit selected rule |
| `d` | Delete selected rule |
| `Space` | Toggle start/stop tunnel |
| `r` | Restart tunnel |
| `l` | Open log view for selected tunnel |
| `?` | Toggle help overlay |
| `q` | Quit |
| `Esc` | Back / cancel |

## Configuration

Rules are persisted automatically. By default, configuration is stored in `~/.config/dwaar/rules.yaml`. You can also manually edit this file to import or modify rules in bulk.

## Development

The project is structured following standard Go project layouts:
- `cmd/dwaar/main.go` - The entrypoint for the application.
- `internal/tunnel` - Tunnel management and SSH forwarding logic.
- `internal/config` - Configuration handling.
- `internal/tui` - Bubble Tea views, components, and styling.

## License

[MIT License](https://opensource.org/licenses/MIT)
