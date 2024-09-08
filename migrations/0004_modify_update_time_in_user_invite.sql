-- +goose Up
ALTER TABLE user_invite
ALTER COLUMN update_time SET DEFAULT NOW();


-- +goose Down
ALTER TABLE user_invite
DROP COLUMN updated_at;