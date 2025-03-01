package bootstrap

import (
	"context"
	"fmt"
	"os"

	"github.com/Hashira21/currency-rate/internal/models/config"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

func DbConnInit(cfg config.Postgres, logger zerolog.Logger) *pgx.Conn {
	user := os.Getenv(cfg.User)
	if user == "" {
		logger.Fatal().Msg("set env variable for database user")
	}

	password := os.Getenv(cfg.Password)
	if password == "" {
		logger.Fatal().Msg("set env variable for database password")
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", user, password, cfg.Host, cfg.Port, cfg.Database)
	pgConnConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	conn, err := pgx.ConnectConfig(context.Background(), pgConnConfig)
	if err != nil {
		logger.Fatal().Msg(err.Error())
		return nil
	}

	err = conn.Ping(context.Background())
	if err != nil {
		logger.Fatal().Msg(err.Error())
		return nil
	}

	logger.Debug().Msg("database connected successfully")

	return conn
}
