package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/price"
)

func (c *Conf) CreateProductStripe(ctx context.Context, amount uint64, productId, productName string) error {
	// Step 1: Retrieve the Stripe secret key from the environment variables
	sKey := os.Getenv("STRIPE_TEST_KEY")
	if sKey == "" {
		// If the key is not set, return an error
		return fmt.Errorf("STRIPE_TEST_KEY not set")
	}

	// Step 2: Assign the Stripe API key to the Stripe library's internal configuration
	stripe.Key = sKey

	// Step 3: Begin a database transaction using the `withTx` method (assumed to be defined elsewhere)
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		// Step 4: Define a SQL query to check if the product already has a Stripe product ID in the database
		sqlQuery := `
				SELECT stripe_product_id 
				FROM product_pricing_stripe
				WHERE product_id = $1
				`

		// Step 5: Declare a variable to hold the Stripe product ID we fetch from the database
		var stripeProductId string

		// Step 6: Execute the query to get the Stripe product ID for the given product ID
		err := tx.QueryRowContext(ctx, sqlQuery, productId).Scan(&stripeProductId)
		if err != nil {
			// Step 7: Handle the case where no rows are found (i.e., the product doesn't have a Stripe product ID yet)
			if !errors.Is(err, sql.ErrNoRows) {
				// If the error is not `sql.ErrNoRows`, that means something went wrong with the query execution; return an error
				return fmt.Errorf("failed to fetch Stripe product ID: %w", err)
			}

			// Step 8: If the product doesn't have a Stripe product ID, create a new Stripe product
			params := &stripe.PriceParams{
				Currency:    stripe.String(string(stripe.CurrencyINR)),
				UnitAmount:  stripe.Int64(int64(amount)),
				Recurring:   &stripe.PriceRecurringParams{}, //this is used for charging customer on recurring basis
				ProductData: &stripe.PriceProductDataParams{Name: stripe.String(productName)},
			}

			// Step 9: Call the Stripe API to create a new product using the parameters
			result, err := price.New(params)
			if err != nil {
				// Log the error and return it if the creation of the Stripe product fails
				slog.Error("failed to create Stripe product", slog.Any("Error", err))
				return fmt.Errorf("failed to create Stripe product: %w", err)
			}

			// Step 10: Define the SQL query to insert the new Stripe product into the `products_stripe` table
			query := `
		INSERT INTO product_pricing_stripe (product_id, stripe_product_id, price_id, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
			// Step 11: Get the current timestamp for record creation and updates
			createdAt := time.Now().UTC()
			updatedAt := createdAt

			// Step 12: Execute the query to insert the new Stripe product into the database
			res, err := tx.ExecContext(ctx, query, productId, result.Product.ID, result.ID, result.UnitAmount, createdAt, updatedAt)
			if err != nil {
				// Log the error and return it if the database insertion fails
				slog.Error("failed to insert Stripe product ID", slog.Any("Error", err))
				return fmt.Errorf("failed to insert Stripe product ID: %w", err)
			}

			// Step 13: Check if the insertion affected any rows (it should affect exactly one row if successful)
			if num, err := res.RowsAffected(); num == 0 || err != nil {
				// If no rows were affected or another error occurred, return an error
				return fmt.Errorf("failed to insert Stripe product ID: %w", err)
			}

			// Step 14: Return `nil` if the Stripe product is successfully added to the database
			return nil
		}

		// Step 15: CustomerId already exist on stripe, no need to add the product
		return nil
	})

	// Step 16: Handle any errors from the transaction function
	if err != nil {
		return err
	}

	// Step 17: If everything succeeds, return `nil`
	return nil
}

