package main

import (
	"fmt"
	"net/http"

	"github.com/Hashira21/currency-rate/internal/bootstrap"
	"github.com/Hashira21/currency-rate/internal/controller"
	"github.com/Hashira21/currency-rate/internal/infrastructure/tech"
	"github.com/Hashira21/currency-rate/internal/providers/frankfurter"
	"github.com/Hashira21/currency-rate/internal/repository/postgres"
	"github.com/Hashira21/currency-rate/internal/router"
	"github.com/Hashira21/currency-rate/internal/service"

	"github.com/gorilla/handlers"
)

var (
	logger = bootstrap.InitLogger()
	cfg    = bootstrap.InitConfig(logger)

	dbConn = bootstrap.DbConnInit(cfg.Postgres, logger)

	frankfurterPrv = frankfurter.NewProvider(&cfg.FrankfurterClient, logger)
	validIsoCodes  = bootstrap.GetValidIsoCodes(frankfurterPrv, logger)
	db             = postgres.New(dbConn, logger)

	svc = service.New(frankfurterPrv, db, logger)
	ctr = controller.New(svc, validIsoCodes, logger)
)

func init() {
	tech.New().SetAppInfo(cfg.Application.Name, cfg.Application.Version)
	bootstrap.StartSyncRates(cfg.SyncRates, svc, logger)
}

func main() {
	// Создаём роутер
	r := router.NewRouter(ctr)

	// Добавляем CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                                // Разрешает запросы с любого домена
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}), // Разрешённые HTTP-методы
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	go svc.AutoUpdateRates() // Запускаем горутину для автообновления курсов
	s := http.Server{
		Addr:         cfg.Application.Port,
		Handler:      corsHandler(r), // Оборачиваем роутер в CORS
		ReadTimeout:  cfg.Application.HttpTimeout,
		WriteTimeout: cfg.Application.HttpTimeout,
	}

	logger.Debug().Msg(fmt.Sprintf("server started on port %s", cfg.Application.Port))

	logger.Fatal().Msg(s.ListenAndServe().Error())
}
