package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AdPatr3s/pre-trainee-practice/internal/store"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Place(ctx context.Context, customerID int, items []store.ItemInput) (int, error) {
	if len(items) == 0 {
		return 0, store.ErrInvalid
	}
	for _, i := range items {
		if i.Qty < 1 {
			return 0, store.ErrInvalid
		}
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var orderID int
	err = tx.QueryRowContext(ctx, "INSERT INTO orders (customer_id, created_at) VALUES ($1, now()) RETURNING id", customerID).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	for _, item := range items {
		var stock int
		var price int
		err = tx.QueryRowContext(ctx, "SELECT stock, price_cents FROM books WHERE id = $1 FOR UPDATE", item.BookID).Scan(&stock, &price)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows){
				return 0, store.ErrNotFound
			}
			return 0, err
		}
		if stock < item.Qty {
			return 0, store.ErrUnavailable
		}
		_, err = tx.ExecContext(ctx, "INSERT INTO order_items (order_id, book_id, qty, unit_price_cents) VALUES ($1, $2, $3, $4)", orderID, item.BookID, item.Qty, price)
		if err != nil {
			return 0, err
		}
		_, err = tx.ExecContext(ctx, "UPDATE books SET stock = stock - $1 WHERE id = $2", item.Qty, item.BookID) 
		if err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}		
	return orderID, nil
}