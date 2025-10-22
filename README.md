# README placeholder
url-shortener/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── api/
│   │   └── handlers.go
│   │
│   ├── id/
│   │   └── idgen.go
│   │
│   ├── util/
│   │   └── base62.go
│   │
│   ├── store/
│   │   ├── store.go
│   │   ├── inmemory.go
│   │   ├── redis_store.go
│   │   ├── cassandra_store.go
│   │   └── postgres_store.go
│   │
│   └── config/
│       └── config.go
│
├── build/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── Makefile
│
├── k8s/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── configmap.yaml
│   ├── hpa.yaml
│   └── statefulset-cassandra.yaml
│
├── scripts/
│   ├── load_test_k6.js
│   ├── init_cassandra.cql
│   └── init_postgres.sql
│
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
└── README.md