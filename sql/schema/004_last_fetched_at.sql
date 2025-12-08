-- +goose Up
ALTER TABLE Feeds ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE Feeds DROP COLUMN last_fetched_at;
