package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RunAddr     string `env:"RUN_ADDRESS"`
	DatabaseURI string `env:"DATABASE_URI"`
	AccSysAddr  string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

// ParseConfig парсит флаги и переменные окружения
func ParseConfig() *Config {
	var cfg Config
	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "address to connect to the database")
	flag.StringVar(&cfg.AccSysAddr, "r", "", "accrual system address")
	flag.Parse()
	cleanenv.ReadEnv(&cfg)
	return &cfg
}
