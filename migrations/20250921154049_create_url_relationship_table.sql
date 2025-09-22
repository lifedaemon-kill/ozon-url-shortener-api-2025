-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS relations
(
    id            SERIAL PRIMARY KEY,
    origin_url    text unique NOT NULL,
    shortened_url text unique NOT NULL,
    created_at    timestamptz DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists relations;
-- +goose StatementEnd
