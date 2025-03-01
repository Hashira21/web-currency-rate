package models

import (
	"database/sql"
	"fmt"
	"time"
)

type CurrencyRate struct {
	Id       string  `db:"id"`
	Currency string  `db:"currency"`
	Base     string  `db:"base"`
	Rate     float32 `db:"rate"`
}

type CurrencyRateDto struct {
	Id       sql.NullString  `db:"id"`
	Currency sql.NullString  `db:"currency"`
	Base     sql.NullString  `db:"base"`
	Rate     sql.NullFloat64 `db:"rate"`
}

type CurrencyRateWithDtDto struct {
	Id       sql.NullString  `db:"id"`
	Currency sql.NullString  `db:"currency"`
	Base     sql.NullString  `db:"base"`
	Rate     sql.NullFloat64 `db:"rate"`
	UpdateDt sql.NullTime    `db:"date"`
}

func (rate *CurrencyRateWithDtDto) FromDto() (CurrencyRateWithDt, error) {
	if !rate.Id.Valid || !rate.Currency.Valid || !rate.Base.Valid || !rate.Rate.Valid || !rate.UpdateDt.Valid {
		return CurrencyRateWithDt{}, fmt.Errorf("params can't be nil %+v", rate)
	}

	return CurrencyRateWithDt{
		Id:       rate.Id.String,
		Currency: rate.Currency.String,
		Base:     rate.Base.String,
		Rate:     rate.Rate.Float64,
		UpdateDt: rate.UpdateDt.Time,
	}, nil
}

func (rate *CurrencyRateWithDtDto) FromDtoToLast() (CurrencyRateLast, error) {
	if !rate.Currency.Valid || !rate.Base.Valid || !rate.Rate.Valid || !rate.UpdateDt.Valid {
		return CurrencyRateLast{}, fmt.Errorf("params can't be nil %+v", rate)
	}

	return CurrencyRateLast{
		Currency: rate.Currency.String,
		Base:     rate.Base.String,
		Rate:     rate.Rate.Float64,
		UpdateDt: rate.UpdateDt.Time,
	}, nil
}

type CurrencyRateWithDt struct {
	Id       string    `db:"id" json:"id" example:"ed7f018b-dc91-4940-8d57-4f91cfe5a8bc"`
	Currency string    `db:"currency" json:"currency" example:"EUR"`
	Base     string    `db:"base" json:"base" example:"USD"`
	Rate     float64   `db:"rate" json:"rate" example:"0.91853"`
	UpdateDt time.Time `db:"date" json:"updateDt" example:"2024-01-20 15:42:12.383064"`
}

type CurrencyRateLast struct {
	Currency  string    `db:"currency" json:"currency" example:"EUR"`
	Base      string    `db:"base" json:"base" example:"USD"`
	Rate      float64   `db:"rate" json:"rate" example:"0.91853"`
	UpdateDt  time.Time `db:"date" json:"updateDt" example:"2024-01-20 15:42:12.383064"`
	ChangePct float64   `json:"changePct" example:"1.23"`
}

type CurrencyRateWithChange struct {
	Id        string    `json:"id"`
	Currency  string    `json:"currency"`
	Base      string    `json:"base"`
	Rate      float64   `json:"rate"`
	UpdateDt  time.Time `json:"updateDt"`
	ChangePct float64   `json:"changePct"` // Изменение в процентах
}
