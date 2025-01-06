-- +goose Up
-- +goose StatementBegin
ALTER TABLE cart_items ADD COLUMN price_id TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cart_items DROP COLUMN price_id;
-- +goose StatementEnd
