package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AdPatr3s/pre-trainee-practice/internal/store"
)

type BookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) *BookRepo {
	return &BookRepo{db: db}
}

func (r *BookRepo) Get(ctx context.Context, id int) (store.Book, error) {
	var b store.Book
	err := r.db.QueryRowContext(ctx, "SELECT id, author_id, title, year, genre, price_cents, stock FROM books WHERE id = $1", id).
		Scan(&b.ID, &b.AuthorID, &b.Title, &b.Year, &b.Genre, &b.PriceCents, &b.Stock)
	if err != nil {
    	if errors.Is(err, sql.ErrNoRows) {
        	return store.Book{}, store.ErrNotFound
    	}
    	return store.Book{}, err
	}
	return b, nil
}

func (r *BookRepo) List(ctx context.Context, genre string) ([]store.Book, error) {
	books := make([]store.Book, 0)
	query := "SELECT id, author_id, title, year, genre, price_cents, stock FROM books"
	var args []any
	if genre != ""{
		query += " WHERE genre = $1"
		args = append(args, genre)
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b store.Book
		err := rows.Scan(&b.ID, &b.AuthorID, &b.Title, &b.Year, &b.Genre, &b.PriceCents, &b.Stock)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepo) UpdateStock(ctx context.Context, id int, stock int) error {
	result, err := r.db.ExecContext(ctx, "UPDATE books SET stock = $1 WHERE id = $2", stock, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return store.ErrNotFound 
	}

	return nil
}