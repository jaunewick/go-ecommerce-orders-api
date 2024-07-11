package handler

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "time"

    "github.com/google/uuid"

    "github.com/Daniel-Giao/orders-api/model"
    "github.com/Daniel-Giao/orders-api/repository/order"
)

type Order struct{
    Repo *order.RedisRepo
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
    var body struct {
        CustomerID uuid.UUID        `json:"customer_id"`
        LineItems  []model.LineItem `json:"line_items"`
    }

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    now := time.Now().UTC()
    order := model.Order{
        OrderID:     rand.Uint64(),
        CustomerID:  body.CustomerID,
        LineItems:   body.LineItems,
        CreatedAt:   &now,
    }

    err := o.Repo.Insert(r.Context(), order)
    if err != nil {
        fmt.Printf("Error inserting order: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    res, err := json.Marshal(order)
    if err != nil {
        fmt.Printf("Error marshalling order: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write(res)
    w.WriteHeader(http.StatusCreated)
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