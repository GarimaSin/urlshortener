package config

import "os"

type Config struct {
    BaseURL          string
    DBType           string
    RedisAddr        string
    PostgresDSN      string
    CassandraHosts   string
    CassandraKeyspace string
    MachineID        string
    Port             string
}

func LoadFromEnv() Config {
    c := Config{
        BaseURL:           getEnv("BASE_URL", "http://localhost:8080"),
        DBType:            getEnv("DB_TYPE", "inmemory"),
        RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
        PostgresDSN:       getEnv("POSTGRES_DSN", "postgres://shortuser:shortpass@localhost:5432/shortdb?sslmode=disable"),
        CassandraHosts:    getEnv("CASSANDRA_HOSTS", "127.0.0.1"),
        CassandraKeyspace: getEnv("CASSANDRA_KEYSPACE", "urlshort"),
        MachineID:         getEnv("MACHINE_ID", "1"),
        Port:              getEnv("PORT", "8080"),
    }
    return c
}

func getEnv(k, def string) string {
    v := os.Getenv(k)
    if v == "" {
        return def
    }
    return v
}
