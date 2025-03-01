package router

import (
	"github.com/gorilla/mux"
)

func NewRouter(ctr Controller) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	techRouter(router)
	setRoutes(router, ctr)

	return router
}
