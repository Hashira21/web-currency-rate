package frankfurter

import (
	"net/http"
	"time"

	"github.com/Hashira21/currency-rate/internal/infrastructure/requester"
	"github.com/Hashira21/currency-rate/internal/models/config"
	"github.com/rs/zerolog"
)

const (
	prvTimeout = 5 * time.Second
)

type provider struct {
	getRate         requester.Requester
	getCurrencyList requester.Requester
	logger          zerolog.Logger
}

func NewProvider(providerCfg *config.Provider, logger zerolog.Logger) *provider {
	httpClient := http.Client{Timeout: prvTimeout}

	return &provider{
		requester.New(&httpClient, *providerCfg, "GetRate"),
		requester.New(&httpClient, *providerCfg, "GetCurrencyList"),
		logger,
	}
}
