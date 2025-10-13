CREATE TABLE IF NOT EXISTS game_data
(
    game_id    INT       NOT NULL,
    user_id    UUID      NOT NULL,
    data       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (game_id, user_id)
);
