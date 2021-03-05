CREATE TABLE IF NOT EXISTS shops
(
    id
            SERIAL
        PRIMARY
            KEY,
    title
            VARCHAR(256),
    address VARCHAR(256),
    phone   VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS products
(
    id
                SERIAL
        PRIMARY
            KEY,
    code
                INTEGER
        NOT
            NULL
        UNIQUE,
    title
                VARCHAR(256),
    description TEXT,
    cost        INTEGER,
    category    VARCHAR(128)
);

CREATE TABLE IF NOT EXISTS receipts
(
    id
        SERIAL
        PRIMARY
            KEY,
    shop_id
        INTEGER
        NOT
            NULL,
    user_id
        INTEGER
        NOT
            NULL,
    cart_item_number
        INTEGER,
    create_date
        TIMESTAMP
        NOT
            NULL
        DEFAULT
            NOW
                (
                )
);

CREATE TABLE IF NOT EXISTS shops_products
(
    id
        SERIAL
        PRIMARY
            KEY,
    shop_id
        INTEGER
        NOT
            NULL,
    product_id
        INTEGER
        NOT
            NULL,
    quantity
        INTEGER
        NOT
            NULL
);

CREATE TABLE IF NOT EXISTS cart_item
(
    id
        SERIAL
        PRIMARY
            KEY,
    product_id
        INTEGER,
    shop_id
        INTEGER
        NOT
            NULL,
    user_id
        INTEGER
        NOT
            NULL,
    quantity
        INTEGER
        NOT
            NULL,
    index
        INTEGER
        NOT
            NULL
);

CREATE TABLE IF NOT EXISTS user_carts
(
    id
        SERIAL
        PRIMARY
            KEY,
    user_id
        INTEGER,
    shop_id
        INTEGER,
    carts_number
        INTEGER,
    UNIQUE
        (
         user_id,
         shop_id
            )
);

CREATE TABLE IF NOT EXISTS receipts_synchro
(
    id
        SERIAL
        PRIMARY
            KEY,
    receipt_id
        INTEGER
        PRIMARY
            KEY,
    is_synchro
        BOOLEAN
);