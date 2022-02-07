package db

import "fmt"

type PostgresDbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type PostgresDbConfigOption func(cfg *PostgresDbConfig)

func WithHost(host string) PostgresDbConfigOption {
	return func(cfg *PostgresDbConfig) {
		cfg.Host = host
	}
} 

func WithPort(port string) PostgresDbConfigOption {
	return func(cfg *PostgresDbConfig) {
		cfg.Port = port
	}
}
 
func WithUsername(username string) PostgresDbConfigOption {
	return func(cfg *PostgresDbConfig) {
		cfg.Username = username
	}
} 

func WithPassword(password string) PostgresDbConfigOption {
	return func(cfg *PostgresDbConfig) {
		cfg.Password = password
	}
} 

func WithDbName(dbname string) PostgresDbConfigOption {
	return func(cfg *PostgresDbConfig) {
		cfg.DBName = dbname
	}
} 

const _disableSSLModeOption = "disable"
func WithDisabledSSLMode() PostgresDbConfigOption {
	return func(cfg *PostgresDbConfig) {
		cfg.SSLMode = _disableSSLModeOption
	}
}

func defaultPostgresDbConfig() *PostgresDbConfig {
	return &PostgresDbConfig{}
}

func (cfg *PostgresDbConfig) toConnectString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)
}