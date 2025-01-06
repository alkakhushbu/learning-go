-- +goose Up
-- +goose StatementBegin
ALTER TABLE cart_items ADD COLUMN price int NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cart_items DROP COLUMN price;
-- +goose StatementEnd
