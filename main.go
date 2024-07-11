package main

import (
    "context"
    "fmt"

    "github.com/Daniel-Giao/orders-api/application"
)

func main() {
    app := application.NewApp()

    err := app.Start(context.TODO())
    if err != nil {
        fmt.Printf("Error starting app: %v\n", err)
    }
}