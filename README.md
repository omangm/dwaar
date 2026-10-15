# Dwaar

Dwaar is a powerful and intuitive Terminal User Interface (TUI) for managing SSH and local port forwarding rules. Built in Go using the renowned [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework, it brings an Elm-like architecture to terminal applications, ensuring a robust and responsive user experience.

## Overview

Managing port forwarding rules via CLI can become cumbersome, especially with multiple connections and SSH hops. Dwaar simplifies this by providing a unified interface to:
- Configure local port to remote host:port mappings.
- Start and stop SSH tunnels seamlessly.
- View real-time connection logs and tunnel statuses.
- Persist configurations for quick setups on subsequent runs.

## Features

- **TUI Management**: Add, edit, and delete port forwarding rules using a built-in form interface, including SSH hop settings and an autostart toggle.
- **Tunnel Control**: Start and stop individual tunnels with a single keystroke.
- **Status Indicators**: Visual cues for `running`, `stopped`, `connecting`, and `error` states.
- **Live Traffic Stats**: Per-tunnel sent/received byte counters while a tunnel is running.
- **Persistent Configuration**: Rules are saved to and loaded from a YAML configuration file.
- **SSH Support**: Forward ports over SSH with key and password authentication.
- **Themes**: Switch between the default and Catppuccin themes at runtime.

## Tech Stack

- **Core Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles)
- **Forms**: [Huh](https://github.com/charmbracelet/huh)
- **SSH Engine**: `golang.org/x/crypto/ssh`

## Installation

Ensure you have [Go](https://golang.org/doc/install) installed (version 1.26+ recommended).

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
| `/` | Filter rules |
| `n` | New rule |
| `e` | Edit selected rule |
| `d` | Delete selected rule |
| `Space` | Toggle start/stop tunnel |
| `l` | Open log view for selected tunnel |
| `g` | Open global log view |
| `t` | Toggle theme |
| `?` | Toggle help overlay |
| `q` | Quit |
| `Esc` | Back / cancel |

## Configuration

Rules are persisted automatically. By default, configuration is stored in `~/.config/dwaar/rules.yaml` (created with user-only `0600` permissions). You can also manually edit this file to import or modify rules in bulk.

Example:

```yaml
version: 1
rules:
  - id: rule_1
    name: db-direct
    local_port: 5432
    remote_host: 10.0.0.5
    remote_port: 5432
    protocol: tcp
    auto_start: true
  - id: rule_2
    name: redis-via-bastion
    local_port: 6379
    remote_host: 10.0.0.6
    remote_port: 6379
    protocol: tcp
    auto_start: false
    via:
      host: bastion.example.com
      port: 22
      user: deploy
      identity_file: ~/.ssh/id_ed25519
      # password: hunter2  # alternative to identity_file
```

> **Note**: an SSH password stored in the config file is kept in plaintext, protected only by file permissions. Prefer key-based authentication when possible.

## Development

The project is structured following standard Go project layouts:
- `cmd/dwaar/main.go` - The entrypoint for the application.
- `internal/tunnel` - Tunnel management and SSH forwarding logic.
- `internal/config` - Configuration handling.
- `internal/tui` - Bubble Tea views, components, and styling.

Run the test suite:

```bash
go test ./...
```

## Website

A bilingual (English / मराठी) marketing site lives in [`site/`](site/). It is a
dependency-free static page — open `site/index.html` directly, or serve the
folder:

```bash
cd site
python3 -m http.server 8080
# then open http://localhost:8080
```

See [`site/README.md`](site/README.md) for details.

## License

[MIT License](https://opensource.org/licenses/MIT)
