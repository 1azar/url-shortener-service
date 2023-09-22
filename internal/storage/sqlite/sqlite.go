package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"url-shortener/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
	    id INTEGER PRIMARY KEY,
	    alias TEXT NOT NULL UNIQUE,
	    url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s Storage) SaveUrl(alias string, urlToStore string) error {
	const fn = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec(urlToStore, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintCheck {
			return fmt.Errorf("%s: %w", fn, storage.ErrURLExists)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s Storage) GetUrl(alias string) (string, error) {
	const fn = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url from url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prapare statment: %w", fn, err)
	}
	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: execute statment: %w", fn, storage.ErrURLNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statment: %w", fn, err)
	}

	return resUrl, nil
}
