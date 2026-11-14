package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/AdPatr3s/pre-trainee-practice/internal/postgres"
	"github.com/AdPatr3s/pre-trainee-practice/internal/store"
)

func setupTest(t *testing.T) *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
    	t.Skip("DATABASE_URL not set")
	}

	db, err := postgres.NewDB(dsn)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
		return nil
	}
	return db
}

func TestBookGet(t *testing.T) {
	db := setupTest(t)
	defer db.Close()

	repo := postgres.NewBookRepo(db)
	ctx := context.Background()

	book, err := repo.Get(ctx, 4)
	if err != nil {
		t.Fatalf("Get(4) failed: %v", err)
	}
	if book.Title != "Emma" {
		t.Errorf("Title = %s, want Emma", book.Title)
	}
	if book.Stock != 4 {
		t.Errorf("Stock = %d, want 4", book.Stock)
	}
}

func TestBookList(t *testing.T) {
	db := setupTest(t)
	defer db.Close()

	repo := postgres.NewBookRepo(db)
	ctx := context.Background()

	books, err := repo.List(ctx, "mystery")

	if err != nil {
		t.Fatalf("List(mystery) failed: %v", err)
	}
	if len(books) != 2 {
		t.Errorf("Got %d books, want 2", len(books))
	}
}

func TestPlaceSuccess(t *testing.T) {
	db := setupTest(t)
	defer db.Close()

	bookRepo := postgres.NewBookRepo(db)
	orderRepo := postgres.NewOrderRepo(db)
	ctx := context.Background()

	items := []store.ItemInput{{BookID: 4, Qty: 1}}
	orderID, err := orderRepo.Place(ctx, 5, items)
	if err != nil {
		t.Fatalf("Place failed: %v", err)
	}
	if orderID <= 0 {
		t.Errorf("Expected orderID > 0, got %d", orderID)
	}
	book, _ := bookRepo.Get(ctx, 4)
	if book.Stock != 3 {
		t.Errorf("Stock = %d, want 3", book.Stock)
	}
	
	bookRepo.UpdateStock(ctx, 4, 4)
	db.Exec("DELETE FROM order_items WHERE order_id = $1", orderID)
	db.Exec("DELETE FROM orders WHERE id = $1", orderID)
}

func TestPlaceUnavailable(t *testing.T) {
	db := setupTest(t)
	defer db.Close()

	orderRepo := postgres.NewOrderRepo(db)
	ctx := context.Background()

	items := []store.ItemInput{{BookID: 4, Qty: 99}}
	_, err := orderRepo.Place(ctx, 5, items)
	
	if !errors.Is(err, store.ErrUnavailable) {
		t.Errorf("Expected ErrUnavailable, got %v", err)
	}
	
	bookRepo := postgres.NewBookRepo(db)
	book, _ := bookRepo.Get(ctx, 4)
	if book.Stock != 4 {
		t.Errorf("Stock changed to %d, want 4 (transaction should have rolled back)", book.Stock)
	}
}