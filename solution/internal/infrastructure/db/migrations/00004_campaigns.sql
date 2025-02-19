-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    advertiser_id UUID NOT NULL REFERENCES advertisers(id) ON DELETE CASCADE,
    impressions_limit BIGINT NOT NULL DEFAULT 0,
    clicks_limit BIGINT NOT NULL DEFAULT 0,
    cost_per_impression DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    cost_per_click DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    ad_title VARCHAR NOT NULL,
    ad_text VARCHAR NOT NULL,
    start_date INT NOT NULL,
    end_date INT NOT NULL,
    pic_id VARCHAR
);

CREATE TABLE IF NOT EXISTS campaigns_targeting (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    gender VARCHAR CHECK (gender IN ('MALE', 'FEMALE', 'ALL') OR gender IS NULL),
    age_from INT,
    age_to INT,
    location VARCHAR
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS campaigns_targeting;
DROP TABLE IF EXISTS campaigns;
-- +goose StatementEnd
