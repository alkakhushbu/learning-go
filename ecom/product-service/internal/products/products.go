package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Conf struct {
	db *sql.DB
}

func NewConf(db *sql.DB) (*Conf, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &Conf{db: db}, nil
}

// InsertProduct inserts a new product into the database and returns the created product.
// It takes a context (`ctx`) and a `NewProduct` struct containing product information.
func (c *Conf) InsertProduct(ctx context.Context, newProduct NewProduct) (Product, error) {
	// Generate a unique ID for the new product using a UUID.
	id := uuid.NewString()

	// Get the current UTC time for `createdAt` and `updatedAt` timestamps for the new product.
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()

	// Declare a `Product` struct to hold the data of the inserted product returned by the SQL query.
	var product Product

	// Use a transaction to ensure atomicity of the database operation.
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		// SQL query to insert a new product into the "products" table.
		// The `RETURNING` clause retrieves the inserted product's data after the operation.
		query := `
				INSERT INTO products 
				(id, name, description, price, category, stock, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id, name, description, price, category, stock, created_at, updated_at;
				`
		// Execute the `INSERT` query within the transaction to add the new product.
		// `QueryRowContext` executes the query and scans the resulting row into the `product` struct.
		err := tx.QueryRowContext(ctx, query,
			id, newProduct.Name, newProduct.Description, newProduct.Price, newProduct.Category, newProduct.Stock, createdAt, updatedAt).
			Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			// Return an error if the query execution or scan fails.
			slog.Error("failed to insert product: %w", err)
			return fmt.Errorf("failed to insert product: %w", err)
		}

		// If the query is successful, return nil to indicate no errors.
		return nil
	})

	// If the transaction or insertion fails, return an error.
	if err != nil {
		return Product{}, fmt.Errorf("failed to insert product: %w", err)
	}

	// Return the inserted product's data as a `Product` struct.
	return product, nil
}

func (c *Conf) GetProductInfo(ctx context.Context, productId string) (ProductInfo, error) {
	var productInfo ProductInfo
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		query := `
				SELECT p.id, p.stock, s.price_id, s.price FROM products p, product_pricing_stripe s
				WHERE p.id = s.product_id AND p.id = $1;
				`
		err := tx.QueryRowContext(ctx, query, productId).
			Scan(&productInfo.Id, &productInfo.Stock, &productInfo.PriceId, &productInfo.Price)

		if err != nil {
			// Return an error if the query execution or scan fails.
			slog.Error("failed to query product info: %w", slog.String("Error", err.Error()))
			return fmt.Errorf("failed to query product info: " + productId)
		}

		// If the query is successful, return nil to indicate no errors.
		return nil
	})
	if err != nil {
		return ProductInfo{}, err
	}
	return productInfo, nil
}

func (c *Conf) DecrementStock(ctx context.Context, productId string, quantity int) error {
	// Decrement stock directly in the database
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		// SQL query to decrement stock atomically
		query := `
		UPDATE products
		SET stock = stock - $1, updated_at = $2
		WHERE id = $3 AND stock > 0
		`

		updatedAt := time.Now().UTC() // Current timestamp

		// Execute the query
		result, err := tx.ExecContext(ctx, query, quantity, updatedAt, productId)
		if err != nil {
			return fmt.Errorf("failed to decrement stock: %w", err)
		}

		// Check if any rows were affected
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to check rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("failed to decrement stock")
		}

		// Successfully decremented stock
		return nil
	})

	// Handle errors
	if err != nil {
		return err
	}

	// Successfully completed
	return nil
}

// withTx is a helper function that simplifies the usage of SQL transactions.
// It begins a transaction using the provided context (`ctx`), executes the given function (`fn`),
// and handles commit or rollback based on the success or failure of the function.
func (c *Conf) withTx(ctx context.Context, fn func(*sql.Tx) error) error {
	// Start a new transaction using the context.
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin tx: %w", slog.String("Error", err.Error()))
		return fmt.Errorf("failed to begin tx: %w", err) // Return an error if the transaction cannot be started.
	}

	// Execute the provided function (`fn`) within the transaction.
	if err := fn(tx); err != nil {
		// If the function returns an error, attempt to roll back the transaction.
		er := tx.Rollback()
		if er != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("failed to rollback withTx: %w", slog.String("Error", err.Error()))
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
