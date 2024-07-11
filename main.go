package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"

    "github.com/Daniel-Giao/orders-api/application"
)

func main() {
    app := application.NewApp(application.LoadConfig())

    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    err := app.Start(ctx)
    if err != nil {
        fmt.Printf("Error starting app: %v\n", err)
    }
}