package controller

import (
	"github.com/rs/zerolog"
)

type controller struct {
	service       Service
	validIsoCodes map[string]struct{}
	logger        zerolog.Logger
}

func New(srv Service, validIsoCodes map[string]struct{}, logger zerolog.Logger) *controller {
	return &controller{
		service:       srv,
		validIsoCodes: validIsoCodes,
		logger:        logger,
	}
}
