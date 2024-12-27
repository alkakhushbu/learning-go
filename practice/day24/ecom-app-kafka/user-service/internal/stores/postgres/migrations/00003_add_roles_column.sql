-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN roles TEXT[] NOT NULL DEFAULT ARRAY['USER'];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN roles;
-- +goose StatementEnd
