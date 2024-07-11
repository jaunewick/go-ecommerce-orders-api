package handler

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "time"
    "strconv"

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
    cursorStr := r.URL.Query().Get("cursor")
    if cursorStr == "" {
        cursorStr = "0"
    }

    const decimal = 10
    const bitSize = 64
    cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    const size = 50
    res, err := o.Repo.FindAll(r.Context(), order.FindAllPage{
        Offset: cursor,
        Size:   size,
    })
    if err != nil {
        fmt.Printf("Error finding all orders: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    var response struct {
        Items []model.Order `json:"items"`
        Next  uint64        `json:"next"`
    }
    response.Items = res.Orders
    response.Next = res.Cursor

    data, err := json.Marshal(response)
    if err != nil {
        fmt.Printf("Error marshalling orders: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write(data)
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