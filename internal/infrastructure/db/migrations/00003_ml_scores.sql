-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ml_scores (
    client_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    advertiser_id UUID NOT NULL REFERENCES advertisers(id) ON DELETE CASCADE,
    score INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ml_scores;
-- +goose StatementEnd
