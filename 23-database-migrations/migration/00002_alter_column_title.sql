-- +goose Up
-- +goose StatementBegin
ALTER TABLE movies ALTER COLUMN title TYPE TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE movies ALTER COLUMN title TYPE VARCHAR(255);
-- +goose StatementEnd
