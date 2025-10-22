package store

import (
    "errors"
    "sync"
    "time"
)

var ErrNotFound = errors.New("not found")

type InMemoryStore struct {
    m  map[string]Record
    mu sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
    return &InMemoryStore{m: map[string]Record{}}
}

func (s *InMemoryStore) Put(short string, long string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.m[short] = Record{ShortCode: short, LongURL: long, CreatedAt: time.Now()}
    return nil
}

func (s *InMemoryStore) Get(short string) (string, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if r, ok := s.m[short]; ok {
        return r.LongURL, nil
    }
    return "", ErrNotFound
}
