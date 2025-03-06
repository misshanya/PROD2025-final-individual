-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR NOT NULL,
    age INTEGER NOT NULL CHECK (age >= 0 AND age <= 200),
    location VARCHAR NOT NULL,
    gender VARCHAR NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
