package config

import (
	"flag"
	"os"
)

var (
	RunAddr     string
	DatabaseURI string
	AccSysAddr  string
)

// ParseConfig парсит флаги и переменные окружения
func ParseConfig() {
	flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&DatabaseURI, "d", "", "address to connect to the database")
	flag.StringVar(&AccSysAddr, "r", "", "accrual system address")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		RunAddr = envRunAddr
	}
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		DatabaseURI = envDatabaseURI
	}
	if envAccSysAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccSysAddr != "" {
		AccSysAddr = envAccSysAddr
	}
}
