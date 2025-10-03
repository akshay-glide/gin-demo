package utils

import (
	"fmt"
	"gin-demo/config"
	"strconv"
)

func optionalOrDefault(val *string, def string) string {
	if val != nil {
		return *val
	}
	return def
}

func GetPostgresDSN(cfg *config.ConfigPostgres) string {
	// If DSN is already provided, use it directly
	if cfg.DSN != nil && *cfg.DSN != "" {
		return *cfg.DSN
	}

	host := *cfg.Host
	port := strconv.Itoa(*cfg.Port)
	user := *cfg.User
	password := *cfg.Password
	dbname := *cfg.Database
	sslmode := optionalOrDefault(cfg.SSLMode, "disable")

	// Construct the DSN string
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=5",
		host, port, user, password, dbname, sslmode,
	)
}
