package tech

import "github.com/rs/zerolog"

type tech struct {
	app    *app
	logger zerolog.Logger
}

type app struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Guid         string   `json:"guid"`
	Hostname     string   `json:"hostname"`
	Dependencies []string `json:"dependencies,omitempty"`
}
