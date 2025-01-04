package carts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const (
	OPEN   = "OPEN"
	CLOSED = "CLOSED"
)

var ErrNotEnoughStock = errors.New("not enough products in stock")
var ErrItemNotInCart = errors.New("item not in cart")

type Conf struct {
	db *sql.DB
}

func NewConf(db *sql.DB) (*Conf, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &Conf{db: db}, nil
}

func (c *Conf) GetCart(ctx context.Context, userId string) (Cart, error) {
	query := `
					SELECT id, user_id, status , created_at, updated_at
					FROM carts 
					WHERE user_id = $1 AND status = $2;`
	var cart Cart
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, query, userId, OPEN).
			Scan(&cart.ID, &cart.UserID, &cart.Status, &cart.CreatedAt, &cart.UpdatedAt)
		if err != nil {
			slog.Error("Cart does not exist", slog.Any("Error", err.Error()))
			return err
		}
		return nil
	})
	if err != nil {
		return Cart{}, err
	}
	return cart, nil
}

func (c *Conf) InsertCart(ctx context.Context, userId string) (Cart, error) {
	query := `
				INSERT INTO 
				carts (id, user_id, status, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, user_id, status, created_at, updated_at`
	var cart Cart
	id := uuid.NewString()
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()

	err := c.withTx(ctx, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, query, id, userId, OPEN, createdAt, updatedAt).
			Scan(&cart.ID, &cart.UserID, &cart.Status, &cart.CreatedAt, &cart.UpdatedAt)
		if err != nil {
			slog.Error("Error in cart creation", slog.Any("Error", err.Error()))
			return err
		}
		return nil
	})
	if err != nil {
		return Cart{}, err
	}
	return cart, nil
}

func (c *Conf) AddItemsToCart(ctx context.Context, cartId string, newItems NewCartItem, stock int) error {
	selectQuery := `SELECT id, product_id, quantity, cart_id, created_at, updated_at
					FROM cart_items
					WHERE cart_id = $1 AND product_id = $2`
	updateQuery := `
				INSERT INTO cart_items (id, product_id, quantity, cart_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (product_id, cart_id)
				DO UPDATE SET quantity = cart_items.quantity + $3
				RETURNING id, product_id, quantity, cart_id, created_at, updated_at;
				`
	var cartItem CartItem
	id := uuid.NewString()
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()

	err := c.withTx(ctx, func(tx *sql.Tx) error {
		// fmt.Println("cartId :", cartId, ".... ProductId", newItems.ProductID)
		err := tx.QueryRowContext(ctx, selectQuery, cartId, newItems.ProductID).
			Scan(&cartItem.ID, &cartItem.ProductID, &cartItem.Quantity,
				&cartItem.CartID, &cartItem.CreatedAt, &cartItem.UpdatedAt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			slog.Error(err.Error())
			return err
		}
		//check final quantity in cart. it must not be greater than the available stock
		if err == nil && cartItem.Quantity+newItems.Quantity > stock {
			slog.Error(ErrNotEnoughStock.Error())
			return ErrNotEnoughStock
		}
		err = tx.QueryRowContext(ctx, updateQuery,
			id, newItems.ProductID, newItems.Quantity, cartId, createdAt, updatedAt).
			Scan(&cartItem.ID, &cartItem.ProductID, &cartItem.Quantity,
				&cartItem.CartID, &cartItem.CreatedAt, &cartItem.UpdatedAt)

		if err != nil {
			slog.Error("Error in cart items creation", slog.Any("Error", err.Error()))
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (c *Conf) RemoveItemsFromCart(ctx context.Context, cartId string, newItems NewCartItem) error {
	selectQuery := `SELECT id, product_id, quantity, cart_id, created_at, updated_at
					FROM cart_items
					WHERE cart_id = $1 AND product_id = $2;`

	updateQuery := `UPDATE cart_items
						SET quantity = quantity - $1, created_at = $2, updated_at = $3
						WHERE id = $4
						RETURNING id, product_id, quantity, cart_id, created_at, updated_at;`

	deleteQuery := `DELETE FROM cart_items
					WHERE quantity <= 0 AND id = $1
					RETURNING id;
				`

	var cartItem CartItem
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()
	// fmt.Println("Cart_id: ", cartId, "...........productId: ", newItems.ProductID)
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, selectQuery, cartId, newItems.ProductID).
			Scan(&cartItem.ID, &cartItem.ProductID, &cartItem.Quantity,
				&cartItem.CartID, &cartItem.CreatedAt, &cartItem.UpdatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				slog.Error(ErrItemNotInCart.Error())
				return ErrItemNotInCart
			}
			fmt.Println(err)
			return err
		}

		id := cartItem.ID
		err = tx.QueryRowContext(ctx, updateQuery, newItems.Quantity, createdAt, updatedAt, id).
			Scan(&cartItem.ID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CartID, &cartItem.CreatedAt, &cartItem.UpdatedAt)
		if err != nil {
			slog.Error("Error in removing items from cart", slog.Any("Error", err.Error()))
			return err
		}

		err = tx.QueryRowContext(ctx, deleteQuery, id).Scan(&id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			slog.Error("Error in deleting items from cart", slog.Any("Error", err.Error()))
			return err

		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil

}

// withTx is a helper function that simplifies the usage of SQL transactions.
// It begins a transaction using the provided context (`ctx`), executes the given function (`fn`),
// and handles commit or rollback based on the success or failure of the function.
func (c *Conf) withTx(ctx context.Context, fn func(*sql.Tx) error) error {
	// Start a new transaction using the context.
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err) // Return an error if the transaction cannot be started.
	}

	// Execute the provided function (`fn`) within the transaction.
	if err := fn(tx); err != nil {
		// If the function returns an error, attempt to roll back the transaction.
		er := tx.Rollback()
		if er != nil && !errors.Is(err, sql.ErrTxDone) {
			// If rollback also fails (and it's not because the transaction is already done),
			// return an error indicating the failure to roll back.
			return fmt.Errorf("failed to rollback withTx: %w", err)
		}
		// Return the original error from the function execution.
		return fmt.Errorf("failed to execute withTx: %w", err)
	}

	// If no errors occur, commit the transaction to apply the changes.
	err = tx.Commit()
	if err != nil {
		// Return an error if the transaction commit fails.
		return fmt.Errorf("failed to commit withTx: %w", err)
	}

	// Return nil if the function executes successfully and the transaction is committed.
	return nil
}
