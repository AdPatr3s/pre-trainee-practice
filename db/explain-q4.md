"Sort  (cost=110.21..112.01 rows=720 width=44)"
"  Sort Key: (COALESCE(sum(order_items.qty), '0'::bigint)) DESC, books.id"
"  ->  HashAggregate  (cost=68.84..76.04 rows=720 width=44)"
"        Group Key: books.id"
"        ->  Hash Right Join  (cost=26.20..59.59 rows=1850 width=40)"
"              Hash Cond: (order_items.book_id = books.id)"
"              ->  Seq Scan on order_items  (cost=0.00..28.50 rows=1850 width=8)"
"              ->  Hash  (cost=17.20..17.20 rows=720 width=36)"
"                    ->  Seq Scan on books  (cost=0.00..17.20 rows=720 width=36)"

"Sort  (cost=110.21..112.01 rows=720 width=44)"
"  Sort Key: (COALESCE(sum(order_items.qty), '0'::bigint)) DESC, books.id"
"  ->  HashAggregate  (cost=68.84..76.04 rows=720 width=44)"
"        Group Key: books.id"
"        ->  Hash Right Join  (cost=26.20..59.59 rows=1850 width=40)"
"              Hash Cond: (order_items.book_id = books.id)"
"              ->  Seq Scan on order_items  (cost=0.00..28.50 rows=1850 width=8)"
"              ->  Hash  (cost=17.20..17.20 rows=720 width=36)"
"                    ->  Seq Scan on books  (cost=0.00..17.20 rows=720 width=36)"

In the Q4 query plan, we were looking for the difference between Seq Scan (sequential scan) and Index Scan (index search).
Without an index on order_items(book_id), the planner was forced to use Hash Join, reading both tables in full.
With an index, the database could use Index Scan, but on such a tiny table (only 13 rows), the PostgreSQL planner often still chooses Seq Scan.
This happens because the overhead of reading the index (additional memory pages) outweighs the benefit of using it on small data sets.
The index will become critically important if the table grows to thousands or millions of rows, to avoid a full brute‑force search.