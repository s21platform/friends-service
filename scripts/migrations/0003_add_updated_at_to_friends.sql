-- +goose Up
ALTER TABLE friends
ADD COLUMN updated_at TIMESTAMP DEFAULT NOW();


-- +goose Down
ALTER TABLE friends
DROP COLUMN updated_at;