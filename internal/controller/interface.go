package controller

import (
	"context"

	"github.com/Hashira21/currency-rate/internal/models"
)

type Service interface {
	GetRateFromProvider(ctx context.Context, toIso, fromIso string) (models.UpdateResponse, error)
	GetById(ctx context.Context, id string) (models.CurrencyRateWithDt, error)
	GetLastRate(ctx context.Context, toIso, fromIso string) (models.CurrencyRateLast, error)
	GetAllLastRates(ctx context.Context) ([]models.CurrencyRateLast, error)
	DeleteByPair(ctx context.Context, currency, base string) error
	UpdateRate(ctx context.Context, currency, base string, rate float64) error
	GetHistory(ctx context.Context, currency, base, period string) ([]models.CurrencyRateWithDt, error)
}
