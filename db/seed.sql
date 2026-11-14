INSERT INTO author (id, name) VALUES
(1, 'George Orwell'),
(2, 'Jane Austen'),
(3, 'J.R.R. Tolkien'),
(4, 'Agatha Christie');

SELECT setval(pg_get_serial_sequence('author', 'id'), (SELECT MAX(id) FROM author));

INSERT INTO books (id, author_id, title, year, genre, price_cents, stock) VALUES
(1, 1, '1984', 1949, 'fiction', 1250, 10),
(2, 1, 'Animal Farm', 1945, 'fiction', 900, 5),
(3, 2, 'Pride and Prejudice', 1813, 'fiction', 1100, 8),
(4, 2, 'Emma', 1815, 'fiction', 1000, 4),
(5, 3, 'The Hobbit', 1937, 'fantasy', 1500, 20),
(6, 3, 'The Lord of the Rings', 1954, 'fantasy', 2500, 7),
(7, 4, 'Murder on the Orient Express', 1934, 'mystery', 850, 12),
(8, 4, 'And Then There Were None', 1939, 'mystery', 800, 6);

SELECT setval(pg_get_serial_sequence('books', 'id'), (SELECT MAX(id) FROM books));

INSERT INTO customers (id, name, email) VALUES
(1, 'Alice', 'alice@example.com'),
(2, 'Bob', 'bob@example.com'),
(3, 'Carol', 'carol@example.com'),
(4, 'Dave', 'dave@example.com'),
(5, 'Eve', 'eve@example.com');

SELECT setval(pg_get_serial_sequence('customers', 'id'), (SELECT MAX(id) FROM customers));

INSERT INTO orders (id, customer_id, created_at) VALUES
(1, 1, '2025-01-10 10:00:00+00'),
(2, 2, '2025-01-20 11:00:00+00'),
(3, 1, '2025-02-05 09:00:00+00'),
(4, 3, '2025-02-14 16:00:00+00'),
(5, 4, '2025-03-01 12:00:00+00'),
(6, 5, '2025-03-03 12:00:00+00'),
(7, 2, '2025-03-15 18:00:00+00'),
(8, 1, '2025-04-01 08:00:00+00'),
(9, 3, '2025-04-18 19:00:00+00'),
(10, 4, '2025-04-20 13:00:00+00');

SELECT setval(pg_get_serial_sequence('orders', 'id'), (SELECT MAX(id) FROM orders));

INSERT INTO order_items (order_id, book_id, qty, unit_price_cents) VALUES
(1, 1, 1, 1250),
(1, 5, 2, 1500),
(2, 2, 1, 900),
(3, 6, 1, 2500),
(4, 7, 2, 850),
(4, 1, 1, 1250),
(5, 5, 1, 1500),
(6, 8, 1, 800),
(7, 3, 1, 1100),
(7, 2, 2, 900),
(8, 1, 3, 1250),
(9, 6, 1, 2500),
(10, 7, 1, 850);