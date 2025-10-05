package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	insertUrlQuery = `INSERT INTO url (url, alias) VALUES ($1, $2) RETURNING id;`
	getUrlQuery    = `SELECT url FROM url WHERE alias = $1;`
	deleteUrlQuery = `DELETE FROM url WHERE alias = $1`
)

type PostgresRepository struct {
	Pool *pgxpool.Pool
}

func (s *PostgresRepository) SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error) {
	var id int64
	err := s.Pool.QueryRow(ctx, insertUrlQuery, urlToSave, alias).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert url")
	}
	return id, nil
}

func (s *PostgresRepository) GetURL(ctx context.Context, alias string) (string, error) {
	var urlStr string // ← просто строка
	err := s.Pool.QueryRow(ctx, getUrlQuery, alias).Scan(&urlStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("alias not found")
		}
		return "", errors.Wrap(err, "failed to get url")
	}
	return urlStr, nil
}

func (s *PostgresRepository) DeleteUrl(ctx context.Context, alias string) error {
	_, err := s.Pool.Exec(ctx, deleteUrlQuery, alias)
	if err != nil {
		return errors.Wrap(err, "failed to delete url")
	}
	return nil
}
