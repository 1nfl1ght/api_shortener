package pgsql

import (
	"api-shorter/internal/storage"
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	connStr := "user=postgres password=Irjkf22013 dbname=api_shortener sslmode=disable"
	const op = "storage.pq.NewStorage"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.pq.SaveURL"

	// stmt, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id`)
	var id int64
	query := "INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id"

	err := s.db.QueryRow(query, urlToSave, alias).Scan(&id)

	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code == "23505" {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
			}
		}
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.pq.GetURL"

	var resUrl string
	err := s.db.QueryRow("select url from url u where u.alias = $1", alias).Scan(&resUrl)

	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			log.Println(pqerr.Code)
			if pqerr.Code == "42601" {
				return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
			}
		}
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resUrl, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.pq.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url u WHERE u.alias = $1")

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: no rows affected", op)
	}

	fmt.Printf("url with alias %s deleted", alias)

	return nil
}
