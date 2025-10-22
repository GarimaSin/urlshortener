package store

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    _ "github.com/lib/pq"
)

type PostgresStore struct {
    db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("open db: %w", err)
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS urls (short_code TEXT PRIMARY KEY, long_url TEXT NOT NULL, created_at TIMESTAMP DEFAULT NOW())`)
    if err != nil {
        return nil, fmt.Errorf("create table: %w", err)
    }
    return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Put(short string, long string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    _, err := s.db.ExecContext(ctx, `INSERT INTO urls (short_code, long_url) VALUES ($1, $2) ON CONFLICT (short_code) DO NOTHING`, short, long)
    return err
}

func (s *PostgresStore) Get(short string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    var long string
    err := s.db.QueryRowContext(ctx, `SELECT long_url FROM urls WHERE short_code = $1 LIMIT 1`, short).Scan(&long)
    if errors.Is(err, sql.ErrNoRows) {
        return "", ErrNotFound
    }
    return long, err
}

func (s *PostgresStore) Close() {
    if s.db != nil {
        _ = s.db.Close()
    }
}
