package store

import (
	"context"
	"errors"
)

var (
    ErrNotFound    = errors.New("store: not found")
    ErrInvalid     = errors.New("store: invalid")
    ErrUnavailable = errors.New("store: unavailable") // not enough stock
)

type Book struct {
	ID         int
	AuthorID   int
	Title      string
	Year       int
	Genre      string
	PriceCents int
	Stock      int
}

type ItemInput struct {
	BookID int
	Qty    int
}

type BookRepository interface {
	Get(ctx context.Context, id int) (Book, error)
	List(ctx context.Context, genre string) ([]Book, error) // empty genre = all
	UpdateStock(ctx context.Context, id int, stock int) error
}

type OrderRepository interface {
	Place(ctx context.Context, customerID int, items []ItemInput) (orderID int, err error)
}