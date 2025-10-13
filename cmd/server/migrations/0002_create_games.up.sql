CREATE TABLE IF NOT EXISTS games
(
    game_id     INT          NOT NULL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    cost        INT          NOT NULL DEFAULT 0,
    category    VARCHAR(100) NOT NULL DEFAULT '',
    home_url    VARCHAR(255) NOT NULL,
    api_key     VARCHAR(255) NOT NULL,
    title_stub  VARCHAR(50) NOT NULL DEFAULT '',
    created_at  TIMESTAMP    NOT NULL,
    updated_at  TIMESTAMP    NOT NULL,

    UNIQUE (title),
    UNIQUE (api_key)
);

INSERT INTO games (game_id, title, description, cost, category, home_url, api_key, title_stub, created_at, updated_at)
VALUES (111111, 'Test Game', 'This is a test game!', 20, 'shooter', 'https://localhost:4443','testkey123', 'test-game',NOW(), NOW());