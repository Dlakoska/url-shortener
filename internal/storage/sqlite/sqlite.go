package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

type SqlLiteRepository struct {
	Db *sql.DB
}

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrUrlExists   = errors.New("url already exists")
	ErrURLNotFound = errors.New("url not found")
)

func (s *SqlLiteRepository) SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveUrl"

	stmt, err := s.Db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, ErrUrlExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *SqlLiteRepository) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"

	stmt, err := s.Db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUrlNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

func (s *SqlLiteRepository) DeleteUrl(ctx context.Context, alias string) error {
	const op = "storage.sqlite.DeleteUrl"

	stmt, err := s.Db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
