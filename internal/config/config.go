package config

import (
	"time"
)

const EnvPath = "local.env"

type Config struct {
	Repository Repository
	LogLevel   string
	Rest       Rest
}
type Rest struct {
	ListenAddress string        `envconfig:"PORT" required:"true"`
	WriteTimeout  time.Duration `envconfig:"WRITE_TIMEOUT" required:"true"`
	ServerName    string        `envconfig:"SERVER_NAME" required:"true"`
	Token         string        `envconfig:"TOKEN" required:"true"`
}

type Postgres struct {
	Host                string        `envconfig:"DB_HOST" required:"true"`
	Port                int           `envconfig:"DB_PORT" required:"true"`
	Name                string        `envconfig:"DB_NAME" required:"true"`
	User                string        `envconfig:"DB_USER" required:"true"`
	Password            string        `envconfig:"DB_PASSWORD" required:"true"`
	SSLMode             string        `envconfig:"DB_SSL_MODE" default:"disable"`
	PoolMaxConns        int           `envconfig:"DB_POOL_MAX_CONNS" default:"5"`
	PoolMaxConnLifetime time.Duration `envconfig:"DB_POOL_MAX_CONN_LIFETIME" default:"180s"`
	PoolMaxConnIdleTime time.Duration `envconfig:"DB_POOL_MAX_CONN_IDLE_TIME" default:"100s"`
}

type SQLite struct {
	Path string `envconfig:"DB_PATH" required:"true"`
}

type Repository struct {
	DbChoice string `envconfig:"DB_CHOICE" required:"true"`
	Postgres Postgres
	SQLite   SQLite
}
