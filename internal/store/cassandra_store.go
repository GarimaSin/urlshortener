package store

import (
    "context"
    "fmt"
    "time"

    "github.com/gocql/gocql"
)

type CassandraStore struct {
    session *gocql.Session
}

func NewCassandraStore(hosts []string, keyspace string) (*CassandraStore, error) {
    cluster := gocql.NewCluster(hosts...)
    cluster.Keyspace = keyspace
    cluster.Consistency = gocql.Quorum
    cluster.Timeout = 5 * time.Second

    session, err := cluster.CreateSession()
    if err != nil {
        return nil, fmt.Errorf("cassandra connect: %w", err)
    }

    // ensure table exists
    err = session.Query(`CREATE TABLE IF NOT EXISTS urls (short_code text PRIMARY KEY, long_url text, created_at timestamp)`).Exec()
    if err != nil {
        return nil, fmt.Errorf("create table: %w", err)
    }
    return &CassandraStore{session: session}, nil
}

func (s *CassandraStore) Put(short string, long string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
    defer cancel()
    q := `INSERT INTO urls (short_code, long_url, created_at) VALUES (?, ?, toTimestamp(now()))`
    return s.session.Query(q, short, long).WithContext(ctx).Exec()
}

func (s *CassandraStore) Get(short string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
    defer cancel()
    var long string
    q := `SELECT long_url FROM urls WHERE short_code = ? LIMIT 1`
    err := s.session.Query(q, short).WithContext(ctx).Scan(&long)
    if err == gocql.ErrNotFound {
        return "", ErrNotFound
    }
    return long, err
}

func (s *CassandraStore) Close() {
    if s.session != nil {
        s.session.Close()
    }
}
