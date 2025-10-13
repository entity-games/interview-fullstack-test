CREATE TABLE IF NOT EXISTS orders
(
    order_id   UUID          NOT NULL PRIMARY KEY,
    user_id    UUID          NOT NULL,
    game_id    INT           NOT NULL DEFAULT 0,
    item_id    VARCHAR(100)  NOT NULL,
    status     SMALLINT      NOT NULL DEFAULT 1,
    type       VARCHAR(50)   NOT NULL,
    content    TEXT,
    created_at TIMESTAMP     NOT NULL,
    updated_at TIMESTAMP     NOT NULL
);

CREATE INDEX idx_orders_user_id ON orders (user_id);
