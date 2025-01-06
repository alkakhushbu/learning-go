-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD COLUMN cart_id UUID;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders DROP COLUMN cart_id;
-- +goose StatementEnd
