package config

import (
	"time"
)

type Config struct {
	Application       Application
	FrankfurterClient Provider
	Postgres          Postgres
	SyncRates         SyncRates
}

type Application struct {
	Name        string
	Version     string
	Port        string
	HttpTimeout time.Duration
}

type Provider struct {
	Host      string
	Endpoints map[string]Endpoint
}

type Endpoint struct {
	Path   string
	Method string
}

type Postgres struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

type SyncRates struct {
	ConfigString string
}
