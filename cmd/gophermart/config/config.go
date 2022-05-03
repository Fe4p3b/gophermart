package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address     string `env:"RUN_ADDRESS,required" envDefault:":8080"`
	AccrualURL  string `env:"ACCRUAL_SYSTEM_ADDRESS,required" envDefault:"http://localhost:8000/api/orders"`
	DatabaseURI string `env:"DATABASE_URI,required" envDefault:"postgres://postgres:12345@localhost:5432/gophermart?sslmode=disable"`
	Secret      string `env:"SECRET" envDefault:"x35k9f"`
}

func SetConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	var (
		address     string
		accrualURL  string
		databaseURI string
	)

	flag.StringVar(&address, "a", "", "Адрес запуска HTTP-сервера")
	flag.StringVar(&accrualURL, "r", "", "Адрес системы расчетов начисления")
	flag.StringVar(&databaseURI, "d", "", "Строка с адресом подключения к БД")
	flag.Parse()

	if address != "" {
		cfg.Address = address
	}

	if accrualURL != "" {
		cfg.AccrualURL = accrualURL
	}

	if databaseURI != "" {
		cfg.DatabaseURI = databaseURI
	}

	return cfg, nil
}
