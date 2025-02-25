package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

// InsertUser inserts a new user into the database and returns the created user.
// It takes a context (`ctx`) and a `NewUser` struct containing user information.
// The function hashes the user's password, inserts the user into the "users" table within a transaction,
// and returns the resulting `User` struct with the inserted data.
func (c *Conf) InsertUser(ctx context.Context, newUser NewUser) (User, error) {
	// Generate a unique ID for the new user using a UUID.
	id := uuid.NewString()

	// Hash the user's password using bcrypt to store it securely in the database.
	// `bcrypt.DefaultCost` determines the cost of the hashing algorithm for computational overhead.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err // Return an error if password hashing fails.
	}

	// Get the current UTC time for `createdAt` and `updatedAt` timestamps for the new user.
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()

	// Declare a `User` struct to hold the data of the inserted user returned by the SQL query.
	var user User

	// Use a transaction to ensure atomicity of the database operation.
	err = c.withTx(ctx, func(tx *sql.Tx) error {
		// SQL query to insert a new user into the "users" table.
		// The `RETURNING` clause retrieves the inserted user's data after the operation.
		query := `
        INSERT INTO users
        (id, name, email, password_hash, created_at, updated_at, roles)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, name, email, created_at, updated_at, roles
        `
		// Execute the `INSERT` query within the transaction to add the new user.
		// `QueryRowContext` executes the query and scans the resulting row into the `user` struct.
		err = tx.QueryRowContext(ctx, query, id, newUser.Name, newUser.Email, hashedPassword, createdAt, updatedAt, newUser.Roles).
			Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Roles)
		if err != nil {
			// Return an error if the query execution or scan fails.
			return fmt.Errorf("failed to insert user: %w", err)
		}

		// If the query is successful, return nil to indicate no errors.
		return nil
	})

	// If the transaction or insertion fails, return an error.
	if err != nil {
		return User{}, fmt.Errorf("failed to insert user: %w", err)
	}

	// Return the inserted user's data as a `User` struct.
	return user, nil
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

func (c *Conf) ValidateUser(ctx context.Context, loginUser LoginUser) (User, error) {
	selectQuery := `select id, name, email, password_hash, created_at, updated_at, roles
					from users 
					where email = $1;`
	var user User
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, selectQuery, loginUser.Email).
			Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.Roles)
		if err != nil {
			slog.Error("user email does not exist", "Error", err.Error())
			return errors.New("wrong email or password")
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginUser.Password))
		if err != nil {
			slog.Error("authentication failed for the user", "Error", err.Error())
			return errors.New("wrong email or password")
		}
		return nil
	})

	if err != nil {
		return User{}, err
	}
	return user, nil
}
