package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	errWrap "github.com/pkg/errors"
	"log/slog"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/postgres"
	"url-shortener/internal/storage/sqlite"
	"url-shortener/pkg/lib/logger/sl"
)

type Repository interface {
	SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error)
	GetURL(ctx context.Context, alias string) (string, error)
	DeleteUrl(ctx context.Context, alias string) error
}

func NewRepository(ctx context.Context, cfg config.Repository) (Repository, error) {
	switch cfg.DbChoice {
	case "postgres":
		repository, err := NewPostgresRepository(ctx, cfg.Postgres)
		if err != nil {
			log.Error("failed to init storage", sl.Err(err))
			return nil, err
		}
		slog.Info("postgres init successfully")
		return repository, nil
	case "sqlite":
		repository, err := NewSQLiteRepository(cfg.SQLite.Path)
		if err != nil {
			log.Error("failed to init storage", sl.Err(err))
			return nil, err
		}
		slog.Info("sqlite init successfully")
		return repository, nil
	}
	log.Error("No DB_CHOICE in config")
	return nil, nil
}

func NewSQLiteRepository(storagePath string) (Repository, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
CREATE TABLE IF NOT EXISTS url(
    id INTEGER PRIMARY KEY, 
    alias TEXT NOT NULL UNIQUE, 
    url TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_alias on url(alias);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &sqlite.SqlLiteRepository{Db: db}, nil
}

func NewPostgresRepository(ctx context.Context, cfg config.Postgres) (Repository, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	// Парсим конфигурацию подключения
	parseConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errWrap.Wrap(err, "failed to parse PostgreSQL config")
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	parseConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, parseConfig)
	if err != nil {
		return nil, errWrap.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &postgres.PostgresRepository{pool}, nil
}
