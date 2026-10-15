package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/omangm/dwaar/internal/config"
	"github.com/omangm/dwaar/internal/tui"
	"github.com/omangm/dwaar/internal/tunnel"
)

func main() {
	path, err := config.DefaultPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving config path: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config %s: %v\n", path, err)
		os.Exit(1)
	}

	mgr := tunnel.NewManager()

	// Start tunnels marked for autostart; failures surface via status events.
	for _, r := range cfg.Rules {
		if r.AutoStart {
			mgr.Start(r)
		}
	}

	p := tea.NewProgram(tui.NewModel(cfg.Rules, mgr), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
