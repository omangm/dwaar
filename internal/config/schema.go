package config

import "github.com/omangm/dwaar/internal/tunnel"

type Config struct {
	Version int                  `yaml:"version"`
	Rules   []tunnel.ForwardRule `yaml:"rules"`
}
