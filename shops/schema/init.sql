CREATE TABLE IF NOT EXISTS shops
(
    id      SERIAL PRIMARY KEY,
    title   VARCHAR(256),
    address VARCHAR(256),
    phone   VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS products
(
    id          SERIAL PRIMARY KEY,
    code        INTEGER NOT NULL UNIQUE,
    title       VARCHAR(256),
    description TEXT,
    cost        INTEGER,
    category    VARCHAR(128)
);

CREATE TABLE IF NOT EXISTS receipts
(
    id          SERIAL PRIMARY KEY,
    cart_id     INTEGER,
    payopt_id   INTEGER,
    create_date TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (cart_id)
);

CREATE TABLE IF NOT EXISTS shops_products
(
    id         SERIAL PRIMARY KEY,
    shop_id    INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity   INTEGER NOT NULL,
    UNIQUE (shop_id, product_id)
);

CREATE TABLE IF NOT EXISTS carts
(
    id      SERIAL PRIMARY KEY,
    shop_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    number  INTEGER NOT NULL,
    UNIQUE (shop_id, user_id, number)
);

CREATE TABLE IF NOT EXISTS cart_item
(
    id         SERIAL PRIMARY KEY,
    product_id INTEGER,
    quantity   INTEGER NOT NULL,
    cart_id    INTEGER,
    UNIQUE (product_id, cart_id)
);

CREATE TABLE IF NOT EXISTS receipts_synchro
(
    id         SERIAL PRIMARY KEY,
    receipt_id INTEGER UNIQUE,
    is_synchro BOOLEAN
);

CREAte TABLE IF NOT EXISTS pay_options
(
    id     SERIAL PRIMARY KEY,
    option VARCHAR(256)
);