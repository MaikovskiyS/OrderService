
CREATE TABLE IF NOT EXISTS orders(
id serial PRIMARY KEY,
order_uid VARCHAR(100),
track_number VARCHAR(100),
entry VARCHAR(100),
locate VARCHAR(100),
internalsignature VARCHAR(100),
customer_id VARCHAR(100),
delivery_service VARCHAR(100),
shardkey VARCHAR(100),
sm_id VARCHAR(100),
date_created VARCHAR(100),
oof_shard VARCHAR(100)
);
CREATE TABLE IF NOT EXISTS delivery(
order_id INT,
name VARCHAR(100),
phone VARCHAR(100),
zip VARCHAR(100),
city VARCHAR(100),
address VARCHAR(100),
region VARCHAR(100),
email VARCHAR(100)
);
CREATE TABLE IF NOT EXISTS payments(
order_id INT,
transaction VARCHAR(100),
request_id VARCHAR(100),
currency VARCHAR(100),
provider VARCHAR(100),
amount REAL,
payment_dt INT,
bank VARCHAR(100),
delivery_cost REAL,
goods_total REAL,
custom_fee REAL
);
CREATE TABLE IF NOT EXISTS items(
order_id INT,
chrt_id REAL,
track_number VARCHAR(100),
price REAL,
rid VARCHAR(100),
name VARCHAR(100),
sale INT,
size VARCHAR(100),
total_price REAL,
nm_id INT,
brand VARCHAR(100),
status INT
);

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS delivery;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS items;
