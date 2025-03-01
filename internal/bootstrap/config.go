package bootstrap

import (
	"github.com/BurntSushi/toml"
	"github.com/Hashira21/currency-rate/internal/models/config"
	"github.com/rs/zerolog"
)

const configPath = "./configs/config.toml"

func InitConfig(logger zerolog.Logger) config.Config {
	var conf config.Config
	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		logger.Fatal().Msg("can`t unmarshall toml configs")
	}

	return conf
}
