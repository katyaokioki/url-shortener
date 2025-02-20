package storage

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls (short VARCHAR(10) PRIMARY KEY, original TEXT UNIQUE);`)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Save(url string) (string, error) {
	var short string
	err := s.db.QueryRow("SELECT short FROM urls WHERE original=$1", url).Scan(&short)
	if err == nil {
		return short, nil
	}
	short, err = generateShortURL()
	if err != nil {
		return "", err
	}
	_, err = s.db.Exec("INSERT INTO urls (short, original) VALUES ($1, $2)", short, url)
	return short, err
}

func (s *PostgresStorage) Get(short string) (string, error) {
	var url string
	err := s.db.QueryRow("SELECT original FROM urls WHERE short=$1", short).Scan(&url)
	if err != nil {
		return "", errors.New("not found")
	}
	return url, nil
}
