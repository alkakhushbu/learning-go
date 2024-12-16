package models

type Book struct {
	ID          int     `json:"id"`    // Unique identifier (primary key in SQL)
	Title       string  `json:"title"` // Title of the book
	AuthorName  string  `json:"author_Name"`
	AuthorEmail string  `json:"author_email"`
	Description string  `json:"description"` // Description of the book
	Category    string  `json:"category"`    // Book category (e.g., Fiction, Biography)
	Price       float64 `json:"price"`       // Price of the book
	Stock       int     `json:"-"`           // Number of copies in stock
}

type NewBook struct {
	Title       string  `json:"title" validate:"required,min=3,max=100"`        // Title is required, length between 3 and 100 characters
	AuthorName  string  `json:"author_name" validate:"required,min=3,max=100"`  // Author's name is required, length between 3 and 100 characters
	AuthorEmail string  `json:"author_email" validate:"required,email"`         // Email is required and must be valid
	Description string  `json:"description" validate:"required,min=10,max=500"` // Description is required, length between 10 and 500 characters
	Category    string  `json:"category" validate:"required,min=3,max=50"`      // Category is required, length between 3 and 50 characters
	Price       float64 `json:"price" validate:"required,gt=0"`                 // Price is required and must be greater than 0
	Stock       int     `json:"stock" validate:"required,gt=0"`                 // Stock is required and must be greater than 0
}

type UpdateBook struct {
	Title       *string  `json:"title" validate:"omitempty,min=3,max=100"`        // Title of the book
	AuthorName  *string  `json:"author_name" validate:"omitempty,min=3,max=100"`  // Author name
	Description *string  `json:"description" validate:"omitempty,min=10,max=500"` // Description
	Category    *string  `json:"category" validate:"omitempty,min=3,max=50"`      // Book category
	Price       *float64 `json:"price" validate:"omitempty,gt=0"`                 // Book price, must be >= 0
	Stock       *int     `json:"stock" validate:"omitempty,gt=0"`                 // Number of copies in stock, must be >= 0
}

// type UpdateBook struct {
// 	Title       *string  `json:"title,omitempty"`       // Title of the book
// 	AuthorName  *string  `json:"author_name,omitempty"` // Author name
// 	Description *string  `json:"description,omitempty"` // Description
// 	Category    *string  `json:"category,omitempty"`    // Book category
// 	Price       *float64 `json:"price,omitempty"`       // Book price, must be >= 0
// 	Stock       *int     `json:"stock,omitempty"`       // Number of copies in stock, must be >= 0
// }
