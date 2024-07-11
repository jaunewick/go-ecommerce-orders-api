package order

import (
    "context"
    "encoding/json"
    "fmt"
    "errors"

    "github.com/redis/go-redis/v9"
    "github.com/Daniel-Giao/orders-api/model"
)

type RedisRepo struct {
    Client *redis.Client
}

func orderIDKey(id uint64) string {
    return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
    data, err := json.Marshal(order)
    if err != nil {
        return fmt.Errorf("error marshalling order: %w", err)
    }

    key := orderIDKey(order.OrderID)

    // Use a transaction to ensure that the order is inserted and added to the set of all orders
    txn := r.Client.TxPipeline()

    res := txn.SetNX(ctx, key, string(data), 0)
    if err := res.Err(); err != nil {
        txn.Discard()
        return fmt.Errorf("error inserting order: %w", err)
    }

    // Add the order ID to the set of all orders
    if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
        txn.Discard()
        return fmt.Errorf("error adding order to set: %w", err)
    }

    if _, err := txn.Exec(ctx); err != nil {
        return fmt.Errorf("error executing transaction: %w", err)
    }

    return nil
}

var ErrNotExist = errors.New("order does not exist")

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
    key := orderIDKey(id)

    value, err := r.Client.Get(ctx, key).Result()
    if errors.Is(err, redis.Nil) {
        return model.Order{}, ErrNotExist
    } else if err != nil {
        return model.Order{}, fmt.Errorf("get order: %w", err)
    }

    var order model.Order
    err = json.Unmarshal([]byte(value), &order)
    if err != nil {
        return model.Order{}, fmt.Errorf("unmarshal order: %w", err)
    }

    return order, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
    key := orderIDKey(id)

    // Use a transaction to ensure that the order is deleted and removed from the set of all orders
    txn := r.Client.TxPipeline()

    err := txn.Del(ctx, key).Err()
    if errors.Is(err, redis.Nil) {
        txn.Discard()
        return ErrNotExist
    } else if err != nil {
        txn.Discard()
        return fmt.Errorf("delete order: %w", err)
    }

    // Remove the order ID from the set of all orders
    if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
        txn.Discard()
        return fmt.Errorf("error removing order from set: %w", err)
    }

    if _, err := txn.Exec(ctx); err != nil {
        return fmt.Errorf("error executing transaction: %w", err)
    }

    return nil
}

func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
    data, err := json.Marshal(order)
    if err != nil {
        return fmt.Errorf("error marshalling order: %w", err)
    }

    key := orderIDKey(order.OrderID)

    err = r.Client.SetXX(ctx, key, string(data), 0).Err()
    if errors.Is(err, redis.Nil) {
        return ErrNotExist
    } else if err != nil {
        return fmt.Errorf("update order: %w", err)
    }

    return nil
}

// FindAllPage represents the pagination parameters for FindAll
type FindAllPage struct {
    Size   uint64
    Offset uint64
}

type FindResult struct {
    Orders []model.Order
    Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
    res := r.Client.SScan(ctx, "orders", page.Offset, "*", int64(page.Size))

    keys, cursor, err := res.Result()
    if err != nil {
        return FindResult{}, fmt.Errorf("error scanning orders: %w", err)
    }

    // If there are no keys, return an empty result
    // MGet will return an error if keys is empty
    if len(keys) == 0 {
        return FindResult{
            Orders: []model.Order{},
        }, nil
    }

    xs, err := r.Client.MGet(ctx, keys...).Result()
    if err != nil {
        return FindResult{}, fmt.Errorf("error getting orders: %w", err)
    }

    orders := make([]model.Order, len(xs))

    for i, x := range xs {
        x := x.(string)
        var order model.Order

        err := json.Unmarshal([]byte(x), &order)
        if err != nil {
            return FindResult{}, fmt.Errorf("unmarshal order: %w", err)
        }

        orders[i] = order
    }

    return FindResult{
        Orders: orders,
        Cursor: cursor,
    }, nil
}