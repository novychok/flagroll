package postgres

import (
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	// Importing the pq package to register the PostgreSQL driver.
	_ "github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Config struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     int    `mapstructure:"POSTGRES_PORT"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	DB       string `mapstructure:"POSTGRES_DB"`
}

type Connection *sql.DB

func New(cfg *Config) (Connection, func(), error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB)

	db, err := otelsql.Open("postgres", dsn)
	if err != nil {
		return nil, func() {}, err
	}

	err = db.Ping()
	if err != nil {
		return nil, func() {}, err
	}

	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		return nil, func() {}, err
	}

	return db, func() {
		_ = db.Close()
	}, nil
}
