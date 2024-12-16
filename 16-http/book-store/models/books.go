package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Conn struct {
	db *pgxpool.Pool
}

func NewConn() (*Conn, error) {
	const (
		host     = "localhost"
		port     = "5433"
		user     = "postgres"
		password = "postgres"
		dbname   = "postgres"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// ParseConfig takes the connection string to generate a config
	config, err := pgxpool.ParseConfig(psqlInfo)
	if err != nil {
		return nil, err
	}

	// MinConns is the minimum number of connections kept open by the pool.
	// The pool will not proactively create this many connections, but once this many have been established,
	// it will not close idle connections unless the total number exceeds MaxConns.
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// MaxConns is the maximum number of connections that can be opened to PostgreSQL.
	// This limit can be used to prevent overwhelming the PostgreSQL server with too many concurrent connections.
	config.MaxConns = 30

	config.HealthCheckPeriod = 5 * time.Minute

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &Conn{db: db}, nil
}

func (c *Conn) InsertBook(ctx context.Context, newBook NewBook) (Book, error) {

	query := `
		INSERT INTO books (
		                   title, author_name,author_email, 
		                   description, category, 
		                   price, stock
		                  
		                   )
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int
	err := c.db.QueryRow(
		ctx, query, newBook.Title, newBook.AuthorName,
		newBook.AuthorEmail, newBook.Description, newBook.Category,
		newBook.Price, newBook.Stock,
	).Scan(&id)

	if err != nil {
		//log.Println(err)
		return Book{}, fmt.Errorf("unable to insert book: %w", err)
	}

	b := Book{
		ID:          id,
		Title:       newBook.Title,
		AuthorName:  newBook.AuthorName,
		AuthorEmail: newBook.AuthorEmail,
		Description: newBook.Description,
		Category:    newBook.Category,
		Price:       newBook.Price,
		Stock:       newBook.Stock,
	}
	return b, nil

}

func (c *Conn) Update(ctx context.Context, id int, updateBook UpdateBook) (Book, error) {

	// Steps
	/*
		1. Add transactions
		2. Add validation to models.Book
		3. If validation fails then rollback the update and report some error to the user
	*/
	selectQuery := `
		SELECT
			id, title, author_name, author_email, description, category, price, stock
		FROM
			books
		WHERE
			id = $1
	`

	var book Book

	// Execute the query and scan the result into the book struct
	err := c.db.QueryRow(ctx, selectQuery, id).Scan(
		&book.ID,
		&book.Title,
		&book.AuthorName,
		&book.AuthorEmail,
		&book.Description,
		&book.Category,
		&book.Price,
		&book.Stock,
	)
	if err != nil {
		return Book{}, fmt.Errorf("unable to fetch book: %w", err)
	}
	data, err := json.Marshal(updateBook)
	if err != nil {
		return Book{}, fmt.Errorf("unable to marshal book: %w", err)
	}
	err = json.Unmarshal(data, &book)
	if err != nil {
		return Book{}, fmt.Errorf("unable to unmarshal book: %w", err)
	}
	//if updateBook.AuthorName != nil {
	//	book.AuthorName = *updateBook.AuthorName
	//}
	//if updateBook.Stock != nil {
	//	book.Stock = *updateBook.Stock
	//}
	//if updateBook.Title != nil {
	//	book.Title = *updateBook.Title
	//}
	//if updateBook.AuthorName != nil {
	//	book.AuthorName = *updateBook.AuthorName
	//}
	//if updateBook.Description != nil {
	//	book.Description = *updateBook.Description
	//}
	//if updateBook.Category != nil {
	//	book.Category = *updateBook.Category
	//}
	//if updateBook.Price != nil {
	//	book.Price = *updateBook.Price
	//}

	query := `
		UPDATE books
		SET title = $1, author_name = $2, description = $3, category = $4, 
		    price = $5, stock = $6
		WHERE id = $7
	`

	// Update the book based on its ID
	_, err = c.db.Exec(ctx, query,
		book.Title, book.AuthorName, book.Description, book.Category,
		book.Price, book.Stock, book.ID,
	)

	if err != nil {
		return Book{}, fmt.Errorf("unable to update book: %w", err)
	}

	fmt.Printf("Book with ID %d updated successfully\n", book.ID)
	return book, nil

}
