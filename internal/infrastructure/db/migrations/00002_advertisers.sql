-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS advertisers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS advertisers;
-- +goose StatementEnd
