    package main

    import (
        "context"
        "log"
        "net/http"
        "os"
        "os/signal"
        "strconv"
        "strings"
        "time"

        "github.com/go-chi/chi/v5"
        "github.com/go-redis/redis/v9"
        "github.com/you/url-shortener/internal/api"
        "github.com/you/url-shortener/internal/config"
        "github.com/you/url-shortener/internal/id"
        "github.com/you/url-shortener/internal/store"
    )

    func main() {
        cfg := config.LoadFromEnv()

        // id generator
        mid := uint16(1)
        if v, err := strconv.Atoi(cfg.MachineID); err == nil {
            mid = uint16(v)
        }
        idgen := id.NewGenerator(mid)

        var backend store.Store
        switch strings.ToLower(cfg.DBType) {
        case "cassandra":
            hosts := []string{cfg.CassandraHosts}
            cass, err := store.NewCassandraStore(hosts, cfg.CassandraKeyspace)
            if err != nil {
                log.Fatalf("cassandra init: %v", err)
            }
            backend = cass
        case "postgres":
            pg, err := store.NewPostgresStore(cfg.PostgresDSN)
            if err != nil {
                log.Fatalf("postgres init: %v", err)
            }
            backend = pg
        default:
            backend = store.NewInMemoryStore()
        }

        // Redis client (used by RedisStore wrapper)
        rcli := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
        cache := store.NewRedisStore(rcli, backend, 24*time.Hour)

        handler := api.NewHandler(cache, idgen, cfg.BaseURL)

        r := chi.NewRouter()
        r.Post("/api/v1/data/shorten", handler.Shorten)
        r.Get("/{shortCode}", handler.Redirect)

        srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

        go func() {
            log.Printf("server starting on :%s
", cfg.Port)
            if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                log.Fatalf("server error: %v", err)
            }
        }()

        quit := make(chan os.Signal, 1)
        signal.Notify(quit, os.Interrupt)
        <-quit
        log.Println("shutting down...")
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        _ = srv.Shutdown(ctx)
    }
