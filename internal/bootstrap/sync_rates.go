package bootstrap

import (
	"github.com/Hashira21/currency-rate/internal/models/config"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type Service interface {
	SyncRates()
}

func StartSyncRates(cfg config.SyncRates, service Service, logger zerolog.Logger) {
	cronJob := cron.New()
	_, err := cronJob.AddFunc(cfg.ConfigString, service.SyncRates)
	if err != nil {
		logger.Error().Msg(err.Error())
	}
	cronJob.Start()
}
