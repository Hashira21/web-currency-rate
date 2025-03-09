package router

import (
	"net/http"

	"github.com/Hashira21/currency-rate/internal/infrastructure/tech"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	apiV1Prefix = "/api/v1"
)

type route struct {
	method  string
	path    string
	name    string
	handler http.HandlerFunc
}

func setRoutes(router *mux.Router, c Controller) {
	var routes = []route{
		{method: http.MethodDelete, path: "/delete/{currency}/{base}", name: "DeleteByPair", handler: c.DeleteByPair},
		{method: http.MethodPut, path: "", name: "UpdateRate", handler: c.UpdateRate},
		{method: http.MethodGet, path: "/by-id/{id}", name: "GetById", handler: c.GetById},
		{method: http.MethodGet, path: "/last", name: "GetLastRate", handler: c.GetLastRate},
		{method: http.MethodGet, path: "/all-last", name: "GetAllLastRates", handler: c.GetAllLastRates},
		{method: http.MethodPatch, path: "/update", name: "UpdateCurrencyRate", handler: c.UpdateCurrencyRate},
		{method: http.MethodGet, path: "/history", name: "GetHistory", handler: c.GetHistory},
	}

	api := router.PathPrefix(apiV1Prefix).Subrouter()

	for _, route := range routes {
		api.
			Name(route.name).
			Methods(route.method).
			Path(route.path).
			Handler(route.handler)
	}
}

func techRouter(router *mux.Router) {
	router.Methods(http.MethodGet).
		Name("prometheus").
		Path("/metrics").
		Handler(promhttp.Handler())

	router.Methods(http.MethodGet).
		Name("GetState").
		Path("/tech/state").
		HandlerFunc(tech.GetState)
}
