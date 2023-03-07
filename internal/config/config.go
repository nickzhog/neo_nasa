package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	Settings struct {
		Address string `env:"SERVER_ADDRESS"`

		PostgresStorage struct {
			DatabaseDSN string `env:"DATABASE_DSN"`
		}

		SourceAPI       string        `env:"SOURCE_API"`
		ApiKey          string        `env:"API_KEY"`
		RequestsIterval time.Duration `env:"REQUESTS_INTERVAL"`
	}
}

func GetConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.Settings.Address, "a", ":8080", "Address for server listen")

	flag.StringVar(&cfg.Settings.PostgresStorage.DatabaseDSN, "d", "", "Database dsn")

	flag.StringVar(&cfg.Settings.SourceAPI, "s", "https://api.nasa.gov/neo/rest/v1/feed", "Source API (NASA)")
	flag.StringVar(&cfg.Settings.ApiKey, "k", "DEMO_KEY", "Key for source API")
	flag.DurationVar(&cfg.Settings.RequestsIterval, "i", time.Millisecond*300, "Interval for NASA update")

	flag.Parse()

	env.Parse(&cfg.Settings)
	env.Parse(&cfg.Settings.PostgresStorage)

	return cfg
}
