# URL Shortener

Scalable URL shortener in Go. Supports three backends: Cassandra, Postgres, Redis (cache) and an in-memory store for local testing.

Endpoints:
- POST /api/v1/data/shorten  (body: {"url": "https://example.com"})
- GET /{shortCode} -> redirect

Quick start (dev):
1. Set environment variables (see .env.example)
2. `go run ./cmd/server`
