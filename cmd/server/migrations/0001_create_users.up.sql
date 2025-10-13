CREATE TABLE IF NOT EXISTS users
(
    user_id    UUID         NOT NULL PRIMARY KEY,
    username   VARCHAR(255) NOT NULL,
    coins      INT          NOT NULL DEFAULT 0,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL,

    UNIQUE (username)
);

INSERT INTO users (user_id, username, coins, created_at, updated_at)
VALUES ('688509d4-c47c-44da-a179-548d6d0dbf50', 'fixture_user', 100, NOW(), NOW());
