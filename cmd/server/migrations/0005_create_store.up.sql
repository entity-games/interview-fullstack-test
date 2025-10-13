CREATE TABLE IF NOT EXISTS store_items
(
    game_id    INT          NOT NULL DEFAULT 0,
    item_id    VARCHAR(200) NOT NULL,
    name       VARCHAR(200) NOT NULL,
    cost       INT                   DEFAULT 0,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL,

    UNIQUE (game_id, item_id)
);

INSERT INTO store_items (game_id, item_id, name, cost, created_at, updated_at)
VALUES (0, 'coins:20', '20x Coins for 10', 10, NOW(), NOW()),
       (111111, 'resource:run', 'Extra runs', 5, NOW(), NOW());
