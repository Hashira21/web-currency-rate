package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Hashira21/currency-rate/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (svc *service) GetRateFromProvider(ctx context.Context, toIso, fromIso string) (models.UpdateResponse, error) {
	respBody, err := svc.frankfurterPrv.GetRate(ctx, toIso, fromIso)
	if err != nil {
		return models.UpdateResponse{}, err
	}

	var rate map[string]interface{}
	if err_ := json.Unmarshal(respBody, &rate); err_ != nil {
		svc.logger.Error().Msg(err_.Error())
		return models.UpdateResponse{}, err_
	}

	if rate["rates"].(map[string]interface{})[toIso] == nil || rate["base"] != fromIso {
		err_ := errors.New("get incorrect value from frankfurter")
		svc.logger.Error().Msg(err_.Error())
		return models.UpdateResponse{}, err_
	}

	currRate := models.CurrencyRate{
		Id:       uuid.New().String(),
		Currency: toIso,
		Base:     fromIso,
		Rate:     float32(rate["rates"].(map[string]interface{})[toIso].(float64)),
	}

	if err_ := svc.db.AddToQueue(ctx, currRate); err_ != nil {
		return models.UpdateResponse{}, err_
	}

	svc.logger.Info().Msg(fmt.Sprintf("succesfully added to queue: %+v", currRate))

	rateResp := models.UpdateResponse{RateId: currRate.Id}

	return rateResp, err
}

func (svc *service) GetById(ctx context.Context, id string) (models.CurrencyRateWithDt, error) {
	return svc.db.GetById(ctx, id)
}

func (svc *service) GetLastRate(ctx context.Context, toIso, fromIso string) (models.CurrencyRateLast, error) {
	return svc.db.GetLastRate(ctx, toIso, fromIso)
}

func (svc *service) SyncRates() {
	ctx := context.Background()

	rate, err := svc.db.ConfirmQueue(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return
		}

		svc.logger.Error().Msg(err.Error())
		return
	}

	svc.logger.Info().Msg(fmt.Sprintf("rate updated successfully: %+v", rate))
}

func (svc *service) GetAllLastRates(ctx context.Context) ([]models.CurrencyRateLast, error) {
	latestRates, err := svc.db.GetAllLastRates(ctx)
	if err != nil {
		return nil, err
	}

	for i := range latestRates {
		prevRate, err := svc.db.GetPreviousRate(ctx, latestRates[i].Currency, latestRates[i].Base)
		if err != nil {
			svc.logger.Warn().Msg(fmt.Sprintf("Не удалось получить предыдущий курс для %s/%s: %v", latestRates[i].Currency, latestRates[i].Base, err))
			continue
		}

		if prevRate.Rate > 0 {
			latestRates[i].ChangePct = ((latestRates[i].Rate - prevRate.Rate) / prevRate.Rate) * 100
		} else {
			latestRates[i].ChangePct = 0
		}
	}

	return latestRates, nil
}

func (svc *service) DeleteByPair(ctx context.Context, currency, base string) error {
	return svc.db.DeleteByPair(ctx, currency, base)
}

func (svc *service) UpdateRate(ctx context.Context, currency, base string, rate float64) error {
	return svc.db.UpdateRate(ctx, currency, base, rate)
}

// Метод для автоматического обновления курсов
func (svc *service) AutoUpdateRates() {
	ticker := time.NewTicker(1 * time.Minute) // Обновление раз в день
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			svc.logger.Info().Msg("Запуск автоматического обновления курсов...")
			err := svc.updateAllRates()
			if err != nil {
				svc.logger.Error().Msg(fmt.Sprintf("Ошибка автообновления курсов: %v", err))
			}
		}
	}
}

// Вспомогательный метод обновления всех курсов
func (svc *service) updateAllRates() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Получаем все доступные валютные пары
	rates, err := svc.db.GetAllLastRates(ctx)
	if err != nil {
		return err
	}

	for _, rate := range rates {
		// Запрашиваем новый курс у API
		respBody, err := svc.frankfurterPrv.GetRate(ctx, rate.Currency, rate.Base)
		if err != nil {
			svc.logger.Warn().Msg(fmt.Sprintf("Не удалось обновить курс для %s/%s: %v", rate.Currency, rate.Base, err))
			continue
		}

		// Разбираем JSON-ответ
		var rateData map[string]interface{}
		if err := json.Unmarshal(respBody, &rateData); err != nil {
			svc.logger.Warn().Msg(fmt.Sprintf("Ошибка парсинга JSON для %s/%s: %v", rate.Currency, rate.Base, err))
			continue
		}

		// Проверяем, есть ли нужные данные в JSON
		if rateData["rates"] == nil || rateData["rates"].(map[string]interface{})[rate.Currency] == nil {
			svc.logger.Warn().Msg(fmt.Sprintf("Некорректный JSON-ответ для %s/%s: %v", rate.Currency, rate.Base, rateData))
			continue
		}

		// Получаем курс
		newRate := rateData["rates"].(map[string]interface{})[rate.Currency].(float64)

		// Сохраняем обновлённый курс в БД
		err = svc.db.UpdateRate(ctx, rate.Currency, rate.Base, newRate)
		if err != nil {
			svc.logger.Warn().Msg(fmt.Sprintf("Ошибка сохранения нового курса %s/%s: %v", rate.Currency, rate.Base, err))
		}
	}

	return nil
}
