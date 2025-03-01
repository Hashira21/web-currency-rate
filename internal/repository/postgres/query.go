package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Hashira21/currency-rate/internal/models"
	"github.com/jackc/pgx/v5"
)

func (db *database) AddToQueue(ctx context.Context, rate models.CurrencyRate) error {
	childCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := db.conn.Exec(childCtx,
		`SELECT * FROM plata_currency_rates.add_to_queue(_id := $1, _currency := $2, _base := $3, _rate := $4)`,
		rate.Id, rate.Currency, rate.Base, rate.Rate)
	if err != nil {
		db.logger.Error().Msg(err.Error())
		return err
	}

	return err
}

func (db *database) ConfirmQueue(ctx context.Context) (models.CurrencyRateWithDt, error) {
	childCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var currRate models.CurrencyRateDto
	var rate models.CurrencyRateWithDtDto

	tx, err := db.conn.Begin(childCtx)
	if err != nil {
		db.logger.Error().Msg(err.Error())
		return models.CurrencyRateWithDt{}, err
	}

	defer tx.Rollback(childCtx)

	err = tx.QueryRow(childCtx,
		`SELECT * FROM plata_currency_rates.confirm_queue();`).
		Scan(&currRate.Id, &currRate.Currency, &currRate.Base, &currRate.Rate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			db.logger.Warn().Msg("queue is empty")
			return models.CurrencyRateWithDt{}, err
		}

		db.logger.Error().Msg(err.Error())
		return models.CurrencyRateWithDt{}, err
	}

	err = tx.QueryRow(childCtx,
		`SELECT * FROM plata_currency_rates.add_to_rates(_id := $1, _currency := $2, _base := $3, _rate := $4);`,
		currRate.Id, currRate.Currency, currRate.Base, currRate.Rate,
	).Scan(&rate.Id, &rate.Currency, &rate.Base, &rate.Rate, &rate.UpdateDt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			db.logger.Warn().Msg(err.Error())
			return models.CurrencyRateWithDt{}, err
		}

		db.logger.Error().Msg(err.Error())
		return models.CurrencyRateWithDt{}, err
	}

	result, err := rate.FromDto()
	if err != nil {
		db.logger.Error().Msg(err.Error())
		return models.CurrencyRateWithDt{}, err
	}

	err = tx.Commit(childCtx)
	if err != nil {
		db.logger.Error().Msg(err.Error())
		return models.CurrencyRateWithDt{}, err
	}

	return result, err
}

func (db *database) GetById(ctx context.Context, id string) (models.CurrencyRateWithDt, error) {
	childCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var rate models.CurrencyRateWithDtDto

	err := db.conn.QueryRow(childCtx,
		`SELECT * FROM plata_currency_rates.get_by_id(_id := $1);`,
		id).
		Scan(&rate.Id, &rate.Currency, &rate.Base, &rate.Rate, &rate.UpdateDt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			db.logger.Warn().Msg(err.Error())
			return models.CurrencyRateWithDt{}, err
		}

		db.logger.Error().Msg(err.Error())
		return models.CurrencyRateWithDt{}, err
	}

	result, err := rate.FromDto()
	if err != nil {
		db.logger.Error().Msg(err.Error())
		return models.CurrencyRateWithDt{}, err
	}

	return result, err
}

func (db *database) GetLastRate(ctx context.Context, toIso, fromIso string) (models.CurrencyRateLast, error) {
	childCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var rate models.CurrencyRateWithDtDto

	err := db.conn.QueryRow(childCtx,
		`SELECT * FROM plata_currency_rates.get_last_rate(_currency := $1, _base := $2);`,
		toIso, fromIso).
		Scan(&rate.Currency, &rate.Base, &rate.Rate, &rate.UpdateDt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			db.logger.Warn().Msg(err.Error())
			return models.CurrencyRateLast{}, err
		}

		db.logger.Error().Msg(err.Error())
		return models.CurrencyRateLast{}, err
	}

	result, err := rate.FromDtoToLast()
	if err != nil {
		db.logger.Error().Msg(err.Error())
		return models.CurrencyRateLast{}, err
	}

	return result, err
}

func (db *database) GetAllLastRates(ctx context.Context) ([]models.CurrencyRateLast, error) {
	childCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	rows, err := db.conn.Query(childCtx,
		`SELECT DISTINCT ON (currency, base) currency, base, rate, date 
		 FROM plata_currency_rates.rates 
		 ORDER BY currency, base, date DESC`)
	if err != nil {
		db.logger.Error().Msg(err.Error())
		return nil, err
	}
	defer rows.Close()

	var rates []models.CurrencyRateLast

	for rows.Next() {
		var rate models.CurrencyRateLast
		err := rows.Scan(&rate.Currency, &rate.Base, &rate.Rate, &rate.UpdateDt)
		if err != nil {
			db.logger.Error().Msg(err.Error())
			return nil, err
		}
		rates = append(rates, rate)
	}

	if err = rows.Err(); err != nil {
		db.logger.Error().Msg(err.Error())
		return nil, err
	}

	return rates, nil
}

func (db *database) DeleteByPair(ctx context.Context, currency, base string) error {
	childCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := db.conn.Exec(childCtx, `
        DELETE FROM plata_currency_rates.rates 
        WHERE currency = $1 AND base = $2
    `, currency, base)

	if err != nil {
		db.logger.Error().Msg(err.Error())
	}

	return err
}

func (db *database) UpdateRate(ctx context.Context, currency, base string, newRate float64) error {
	childCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `
        INSERT INTO plata_currency_rates.rates (id, currency, base, rate, date)
        VALUES (gen_random_uuid(), $1, $2, $3, NOW());
    `

	_, err := db.conn.Exec(childCtx, query, currency, base, newRate)
	if err != nil {
		db.logger.Error().Msg(fmt.Sprintf("Ошибка добавления нового курса %s/%s: %v", currency, base, err))
		return err
	}

	db.logger.Info().Msg(fmt.Sprintf("Новый курс %s/%s успешно добавлен: %f", currency, base, newRate))
	return nil
}

func (db *database) GetLastRateWithChange(ctx context.Context, toIso, fromIso string) (models.CurrencyRateWithChange, error) {
	childCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var rate models.CurrencyRateWithDtDto
	var prevRate models.CurrencyRateWithDtDto

	// Получаем последний курс
	err := db.conn.QueryRow(childCtx,
		`SELECT id, currency, base, rate, date FROM plata_currency_rates.rates
         WHERE currency = $1 AND base = $2 ORDER BY date DESC LIMIT 1;`,
		toIso, fromIso).
		Scan(&rate.Id, &rate.Currency, &rate.Base, &rate.Rate, &rate.UpdateDt)
	if err != nil {
		return models.CurrencyRateWithChange{}, err
	}

	// Получаем предыдущий курс
	err = db.conn.QueryRow(childCtx,
		`SELECT rate FROM plata_currency_rates.rates
         WHERE currency = $1 AND base = $2 ORDER BY date DESC OFFSET 1 LIMIT 1;`,
		toIso, fromIso).
		Scan(&prevRate.Rate)
	if err != nil {
		prevRate.Rate.Float64 = rate.Rate.Float64 // Если предыдущего курса нет, берём тот же
	}

	// Вычисляем процентное изменение
	changePct := ((rate.Rate.Float64 - prevRate.Rate.Float64) / prevRate.Rate.Float64) * 100

	return models.CurrencyRateWithChange{
		Id:        rate.Id.String,
		Currency:  rate.Currency.String,
		Base:      rate.Base.String,
		Rate:      rate.Rate.Float64,
		UpdateDt:  rate.UpdateDt.Time,
		ChangePct: changePct,
	}, nil
}

func (db *database) GetPreviousRate(ctx context.Context, currency, base string) (models.CurrencyRateLast, error) {
	childCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var prevRate models.CurrencyRateWithDtDto

	err := db.conn.QueryRow(childCtx,
		`SELECT id, currency, base, rate, date FROM plata_currency_rates.rates 
         WHERE currency = $1 AND base = $2 
         ORDER BY date DESC OFFSET 1 LIMIT 1;`,
		currency, base).
		Scan(&prevRate.Id, &prevRate.Currency, &prevRate.Base, &prevRate.Rate, &prevRate.UpdateDt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			db.logger.Warn().Msg(fmt.Sprintf("Нет предыдущего курса для %s/%s", currency, base))
			return models.CurrencyRateLast{}, err
		}
		db.logger.Error().Msg(err.Error())
		return models.CurrencyRateLast{}, err
	}

	result, err := prevRate.FromDtoToLast()
	if err != nil {
		db.logger.Error().Msg(err.Error())
		return models.CurrencyRateLast{}, err
	}

	return result, nil
}
