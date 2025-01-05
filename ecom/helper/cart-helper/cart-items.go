package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

type CartItem struct {
    ID        int
    ProductID int
    Quantity  int
    CartID    int
    CreatedAt sql.NullTime
    UpdatedAt sql.NullTime
}

func fetchCartItems(ctx context.Context, db *sql.DB, cartID int) ([]CartItem, error) {
    query := `SELECT id, product_id, quantity, cart_id, created_at, updated_at
                FROM cart_items
                WHERE cart_id = $1`

    rows, err := db.QueryContext(ctx, query, cartID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cartItems []CartItem
    for rows.Next() {
        var item CartItem
        if err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.CartID, &item.CreatedAt, &item.UpdatedAt); err != nil {
            return nil, err
        }
        cartItems = append(cartItems, item)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return cartItems, nil
}

func main() {
    connStr := "user=youruser dbname=yourdb sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ctx := context.Background()
    cartID := 1 // Replace with your cart ID
    cartItems, err := fetchCartItems(ctx, db, cartID)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Fetched %d cart items\n", len(cartItems))
    for _, item := range cartItems {
        fmt.Printf("%+v\n", item)
    }
}
