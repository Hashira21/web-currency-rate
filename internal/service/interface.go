package service

import (
	"context"

	"github.com/Hashira21/currency-rate/internal/models"
)

type FrankfurterPrv interface {
	GetRate(ctx context.Context, toIso, fromIso string) ([]byte, error)
}

type Postgres interface {
	AddToQueue(ctx context.Context, rate models.CurrencyRate) error
	ConfirmQueue(ctx context.Context) (models.CurrencyRateWithDt, error)
	GetById(ctx context.Context, id string) (models.CurrencyRateWithDt, error)
	GetLastRate(ctx context.Context, toIso, fromIso string) (models.CurrencyRateLast, error)
	GetPreviousRate(ctx context.Context, currency, base string) (models.CurrencyRateLast, error)
	GetAllLastRates(ctx context.Context) ([]models.CurrencyRateLast, error)
	DeleteByPair(ctx context.Context, currency, base string) error
	UpdateRate(ctx context.Context, currency, base string, rate float64) error
}
