-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS author (
    id bigserial PRIMARY KEY,
    title varchar(255) NOT NULL,
    year integer NOT NULL,
    runtime integer NOT NULL,
    genres text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS movies;
-- +goose StatementEnd
