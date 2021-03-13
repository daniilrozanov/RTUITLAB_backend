CREATE TABLE IF NOT EXISTS products
(
    id    SERIAL PRIMARY KEY,
    date  TIMESTAMP NOT NULL DEFAULT NOW(),
    title VARCHAR(128),
    cost  INT
);

CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(64) UNIQUE ,
    password_hash VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS users_products
(
    id         SERIAL PRIMARY KEY,
    user_id    INT,
    product_id INT,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS receipts
(
    id    SERIAL PRIMARY KEY,
    body  jsonb
);

CREATE TABLE IF NOT EXISTS users_receipts
(
    id         SERIAL PRIMARY KEY,
    user_id    INT,
    receipt_id INT,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (receipt_id) REFERENCES receipts (id) ON DELETE CASCADE
);
