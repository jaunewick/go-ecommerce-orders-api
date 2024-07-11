package handler

import (
    "fmt"
    "net/http"
)

type Order struct{}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Create an order\n")
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("List all orders\n")
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Get an order by ID\n")
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Update an order by ID\n")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Delete an order by ID\n")
}