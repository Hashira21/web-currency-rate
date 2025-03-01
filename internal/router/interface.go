package router

import "net/http"

type Controller interface {
	UpdateRate(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	GetLastRate(w http.ResponseWriter, r *http.Request)
	GetAllLastRates(w http.ResponseWriter, r *http.Request)
	DeleteByPair(w http.ResponseWriter, r *http.Request)
	UpdateCurrencyRate(w http.ResponseWriter, r *http.Request)
}
