package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Hashira21/currency-rate/internal/infrastructure/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// UpdateRate godoc
// @Summary      	Send signal to update rate
// @Tags         	Methods
// @Param 			rate query string false "currency rate" example(EUR/USD)
// @Success      	200 {object} models.UpdateResponse "success"
// @Failure      	400 "validation error"
// @Failure      	500 "service unavailable"
// @Router       	/ [put]
func (ctr *controller) UpdateRate(w http.ResponseWriter, r *http.Request) {
	currencyRate := r.URL.Query().Get("rate")
	currencies := strings.Split(currencyRate, "/")
	if len(currencies) != 2 {
		err_ := errors.New("parameter doesn't match pattern EUR/USD")
		ctr.logger.Error().Msg(err_.Error())
		response.WriteError(w, http.StatusBadRequest, err_)
		return
	}

	if invalidIso, isInvalid := ctr.validateIsoCode(&currencies[0], &currencies[1]); !isInvalid {
		err_ := fmt.Errorf("uexpected iso code %s. try this one: %s", invalidIso, ctr.getValidIsoCodesString())
		ctr.logger.Error().Msg(fmt.Sprintf("uexpected iso code %s", invalidIso))
		response.WriteError(w, http.StatusBadRequest, err_)
		return
	}

	rateId, err := ctr.service.GetRateFromProvider(r.Context(), currencies[0], currencies[1])
	if err != nil {
		ctr.logger.Error().Msg(err.Error())
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	respBody, err := json.Marshal(rateId)
	if err != nil {
		ctr.logger.Error().Msg(err.Error())
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, respBody)
}

// GetById godoc
// @Summary      	Get currency rate by id
// @Tags         	Methods
// @Param 			id path string false "currency rate update ID" example(ed7f018b-dc91-4940-8d57-4f91cfe5a8bc)
// @Success      	200 {object} models.CurrencyRateWithDt "success"
// @Failure      	400 "validation error"
// @Failure      	500 "service unavailable"
// @Router       	/by-id/{id} [get]
func (ctr *controller) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		err_ := errors.New("set id value")
		ctr.logger.Error().Msg(err_.Error())
		response.WriteError(w, http.StatusBadRequest, err_)
		return
	}

	if err_ := uuid.Validate(id); err_ != nil {
		ctr.logger.Error().Msg(err_.Error())
		response.WriteError(w, http.StatusBadRequest, err_)
		return
	}

	result, err := ctr.service.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		ctr.logger.Error().Msg(err.Error())
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	respBody, err := json.Marshal(result)
	if err != nil {
		ctr.logger.Error().Msg(err.Error())
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, respBody)
}

// GetLastRate godoc
// @Summary      	Get latest currency rate
// @Tags         	Methods
// @Param 			rate query string false "currency rate" example(EUR/USD)
// @Success      	200 {object} models.CurrencyRateLast "success"
// @Failure      	400 "validation error"
// @Failure      	500 "service unavailable"
// @Router       	/last [get]
func (ctr *controller) GetLastRate(w http.ResponseWriter, r *http.Request) {
	currencyRate := r.URL.Query().Get("rate")
	currencies := strings.Split(currencyRate, "/")
	if len(currencies) != 2 {
		err_ := errors.New("parameter doesn't match pattern EUR/USD")
		ctr.logger.Error().Msg(err_.Error())
		response.WriteError(w, http.StatusBadRequest, err_)
		return
	}

	if invalidIso, isInvalid := ctr.validateIsoCode(&currencies[0], &currencies[1]); !isInvalid {
		err_ := fmt.Errorf("uexpected iso code %s. try this one: %s", invalidIso, ctr.getValidIsoCodesString())
		ctr.logger.Error().Msg(fmt.Sprintf("uexpected iso code %s", invalidIso))
		response.WriteError(w, http.StatusBadRequest, err_)
		return
	}

	result, err := ctr.service.GetLastRate(r.Context(), currencies[0], currencies[1])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		ctr.logger.Error().Msg(err.Error())
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	respBody, err := json.Marshal(result)
	if err != nil {
		ctr.logger.Error().Msg(err.Error())
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, respBody)
}

func (ctr *controller) GetAllLastRates(w http.ResponseWriter, r *http.Request) {
	result, err := ctr.service.GetAllLastRates(r.Context())
	if err != nil {
		ctr.logger.Error().Msg(err.Error())
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	respBody, err := json.Marshal(result)
	if err != nil {
		ctr.logger.Error().Msg(err.Error())
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, respBody)
}

func (ctr *controller) DeleteByPair(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	currency := vars["currency"]
	base := vars["base"]

	if currency == "" || base == "" {
		http.Error(w, "Неверные параметры валютной пары", http.StatusBadRequest)
		return
	}

	err := ctr.service.DeleteByPair(r.Context(), currency, base)
	if err != nil {
		ctr.logger.Error().Msg(err.Error())
		http.Error(w, "Ошибка удаления курса", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateRateById godoc
// @Summary       Update currency rate by ID
// @Tags          Methods
// @Param         id path string true "Currency Rate ID"
// @Param         rate body models.UpdateRateRequest true "New currency rate"
// @Success       200 {object} models.UpdateResponse "success"
// @Failure       400 "validation error"
// @Failure       500 "service unavailable"
// @Router        /update/{id} [patch]
func (ctr *controller) UpdateCurrencyRate(w http.ResponseWriter, r *http.Request) {
	currency := r.URL.Query().Get("currency")
	base := r.URL.Query().Get("base")
	rateStr := r.URL.Query().Get("rate")

	if currency == "" || base == "" || rateStr == "" {
		err_ := errors.New("не указаны параметры валюты, базы или курса")
		ctr.logger.Error().Msg(err_.Error())
		response.WriteError(w, http.StatusBadRequest, err_)
		return
	}

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil || rate <= 0 {
		err_ := errors.New("некорректное значение курса")
		ctr.logger.Error().Msg(err_.Error())
		response.WriteError(w, http.StatusBadRequest, err_)
		return
	}

	err = ctr.service.UpdateRate(r.Context(), currency, base, rate)
	if err != nil {
		ctr.logger.Error().Msg(err.Error())
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, []byte(`{"message": "Курс обновлён успешно"}`))
}
