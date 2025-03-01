package service

import (
	"github.com/rs/zerolog"
)

type service struct {
	frankfurterPrv FrankfurterPrv
	db             Postgres
	logger         zerolog.Logger
}

func New(frankfurterPrv FrankfurterPrv, db Postgres, logger zerolog.Logger) *service {
	return &service{
		frankfurterPrv: frankfurterPrv,
		db:             db,
		logger:         logger,
	}
}
