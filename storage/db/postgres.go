package db

import "github.com/jmoiron/sqlx"

const _postgresDriverName = "postgres"
func ConnectToPostgresDb(opts ...PostgresDbConfigOption) (*sqlx.DB, error) {
	cfg := defaultPostgresDbConfig() 
	for _, opt := range opts {
		opt(cfg)
	} 

	db, err := sqlx.Connect(_postgresDriverName, cfg.toConnectString())
	if err != nil {
		return nil, err
	} 

	return db, nil
}