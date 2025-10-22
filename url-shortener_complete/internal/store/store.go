package store

import "time"

type Record struct {
    ShortCode string
    LongURL   string
    CreatedAt time.Time
}

type Store interface {
    Put(short string, long string) error
    Get(short string) (string, error)
}
