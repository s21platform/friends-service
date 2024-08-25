-- +goose Up
CREATE TABLE IF NOT EXISTS user_invite (
                                           id SERIAL PRIMARY KEY,
                                           initiator UUID,
                                           invited VARCHAR(255),
    is_closed BOOLEAN DEFAULT FALSE,
    update_time TIMESTAMP
    );






-- +goose Down
DROP TABLE IF EXISTS user_invite;