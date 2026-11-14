CREATE TABLE author (
id SERIAL PRIMARY KEY,
name TEXT NOT NULL UNIQUE
);
CREATE TABLE books (
	id SERIAL PRIMARY KEY,
	author_id INT NOT NULL REFERENCES author(id),
	title TEXT NOT NULL,
	year INT NOT NULL CHECK (year >= 1000 AND year <= 2100),
	genre TEXT NOT NULL,
    price_cents INT NOT NULL CHECK (price_cents >= 0),
    stock INT NOT NULL CHECK (stock >= 0),
    UNIQUE (author_id, title)
);

CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE order_items (
    order_id INT NOT NULL REFERENCES orders(id),
    book_id INT NOT NULL REFERENCES books(id),
    qty INT NOT NULL CHECK (qty > 0),
    unit_price_cents INT NOT NULL CHECK (unit_price_cents >= 0),
    PRIMARY KEY (order_id, book_id)
);

CREATE INDEX idx_books_genre ON books(genre);
CREATE INDEX idx_books_author_id ON books(author_id);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_order_items_book_id ON order_items(book_id);