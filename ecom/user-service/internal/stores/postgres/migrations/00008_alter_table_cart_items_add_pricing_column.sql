-- +goose Up
-- +goose StatementBegin
ALTER TABLE cart_items ADD COLUMN price int NOT NULL default 10000;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cart_items DROP COLUMN price;
-- +goose StatementEnd
