package bootstrap

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
)

type FrankfurterPrv interface {
	GetCurrencyList(ctx context.Context) ([]byte, error)
}

func GetValidIsoCodes(frankPrv FrankfurterPrv, logger zerolog.Logger) map[string]struct{} {
	ctx := context.Background()

	res, err := frankPrv.GetCurrencyList(ctx)
	if err != nil {
		logger.Fatal().Msg("failed to get a list of valid iso codes")
	}

	var validIso map[string]interface{}
	if err_ := json.Unmarshal(res, &validIso); err_ != nil {
		logger.Fatal().Msg("error in unmarshalling ")
	}

	validIsoCash := make(map[string]struct{}, len(validIso))
	for i := range validIso {
		validIsoCash[i] = struct{}{}
	}

	return validIsoCash
}
