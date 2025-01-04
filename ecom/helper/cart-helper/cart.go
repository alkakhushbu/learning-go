package main

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/mux"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var db *sqlx.DB

func main() {
    var err error
    db, err = sqlx.Connect("postgres", "user=youruser dbname=yourdb sslmode=disable")
    if err != nil {
        panic(err)
    }

    r := mux.NewRouter()
    r.HandleFunc("/carts", createOrUpdateCart).Methods("POST")
    http.ListenAndServe(":8080", r)
}

func createOrUpdateCart(w http.ResponseWriter, r *http.Request) {
    type CartRequest struct {
        ProductID int `json:"product_id"`
        Quantity  int `json:"quantity"`
    }

    var req CartRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    authHeader := r.Header.Get("Authorization")
    tokenString := authHeader[len("Bearer "):]

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("your-secret-key"), nil
    })

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := int(claims["user_id"].(float64))

        var cartID int
        err = db.Get(&cartID, "SELECT id FROM carts WHERE user_id = $1 AND status = $2", userID, "OPEN")
        if err != nil {
            result, err := db.Exec("INSERT INTO carts (user_id, created_at, updated_at, status) VALUES ($1, $2, $3, $4)", userID, time.Now(), time.Now(), "OPEN")
            if err != nil {
                http.Error(w, "Failed to create cart", http.StatusInternalServerError)
                return
            }
            cartID, err = result.LastInsertId()
            if err != nil {
                http.Error(w, "Failed to retrieve cart ID", http.StatusInternalServerError)
                return
            }
        }

        _, err = db.Exec("INSERT INTO cart_items (product_id, quantity, cart_id) VALUES ($1, $2, $3) ON CONFLICT (product_id, cart_id) DO UPDATE SET quantity = cart_items.quantity + $2", req.ProductID, req.Quantity, cartID)
        if err != nil {
            http.Error(w, "Failed to add product to cart", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Product added to cart"})
    } else {
        http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
    }
}
