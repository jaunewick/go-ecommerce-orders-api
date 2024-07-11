package application

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/redis/go-redis/v9"
)

type App struct {
    router http.Handler
    rdb    *redis.Client
    config Config
}

func NewApp(config Config) *App {
    app := &App{
        rdb:        redis.NewClient(&redis.Options{
            Addr:   config.RedisAddr,
        }),
        config: config,
    }

    app.loadRoutes()

    return app
}

func (a *App) Start(ctx context.Context) error {
    server := &http.Server{
        Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
        Handler: a.router,
    }

    err := a.rdb.Ping(ctx).Err()
    if err != nil {
        return fmt.Errorf("error connecting to redis: %w", err)
    }

    // Close the redis connection when the app stops
    defer func() {
        if err := a.rdb.Close(); err != nil {
            fmt.Println("error closing redis connection\n", err)
        }
    }()

    fmt.Printf("Server listening on port %s\n", server.Addr)

    // Create a channel to listen for errors from the server
    ch := make(chan error, 1)

    // Start the server in a goroutine
    go func() {
        err = server.ListenAndServe()
        if err != nil {
            ch <- fmt.Errorf("error starting server: %w", err)
        }
        close(ch)
    }()

    // Listen for context cancellation
    select {
    case err = <-ch:
        return err
    case <-ctx.Done():
        timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        return server.Shutdown(timeout)
    }

    return nil
}