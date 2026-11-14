-- Q1: revenue by month (UTC)
SELECT 
    to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM') AS year_month,
    SUM(qty * unit_price_cents) AS revenue_cents
FROM orders
JOIN order_items ON orders.id = order_items.order_id
GROUP BY year_month
ORDER BY year_month;

--Q2 — top customers by spend
SELECT customers.id AS customer_id,
	customers.name,
	SUM(qty*unit_price_cents) AS spend_cents
FROM customers
JOIN orders ON customers.id = orders.customer_id
JOIN order_items ON orders.id = order_items.order_id
GROUP BY customers.name, customers.id 
ORDER BY spend_cents DESC, customer_id ASC;