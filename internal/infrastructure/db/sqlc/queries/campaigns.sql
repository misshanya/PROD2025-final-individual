-- name: CreateCampaign :one
INSERT INTO campaigns (
    advertiser_id,
    impressions_limit, clicks_limit,
    cost_per_impression, cost_per_click,
    ad_title, ad_text,
    start_date, end_date
) VALUES (
    @advertiser_id::uuid,
    @impressions_limit::bigint, @clicks_limit::bigint,
    @cost_per_impression::decimal(10,2), @cost_per_click::decimal(10,2),
    @ad_title::varchar, @ad_text::varchar,
    @start_date::int, @end_date::int
)
RETURNING *;

-- name: CreateCampaignTargeting :one
INSERT INTO campaigns_targeting (
    campaign_id,
    gender,
    age_from, age_to,
    location
) VALUES (
    @campaign_id::uuid,
    COALESCE(sqlc.narg(gender)::varchar, NULL),
    COALESCE(sqlc.narg(age_from)::int, NULL), COALESCE(sqlc.narg(age_to)::int, NULL),
    COALESCE(sqlc.narg(location)::varchar, NULL)
)
RETURNING *;

-- name: SetCampaignPicture :exec
UPDATE campaigns
SET
    pic_id = @picture_id::varchar
WHERE
    id = @campaign_id::uuid;

-- name: GetCampaignPicID :one
SELECT pic_id FROM campaigns
WHERE id = @campaign_id::uuid; 

-- name: GetCampaignsWithTargetingByAdvertiserID :many
SELECT * FROM campaigns JOIN campaigns_targeting ON campaigns.id = campaigns_targeting.campaign_id
WHERE advertiser_id = @advertiser_id::uuid
LIMIT $1 OFFSET $2;

-- name: GetCampaignWithTargetingByID :one
SELECT * FROM campaigns JOIN campaigns_targeting ON campaigns.id = campaigns_targeting.campaign_id
WHERE campaigns.id = @campaign_id::uuid;

-- name: UpdateCampaign :one
UPDATE campaigns
SET
    impressions_limit = @impressions_limit::bigint, clicks_limit = @clicks_limit::bigint,
    cost_per_impression = @cost_per_impression::decimal(10,2), cost_per_click = @cost_per_click::decimal(10,2),
    ad_title = @ad_title::varchar, ad_text = @ad_text::varchar,
    start_date = @start_date::int, end_date = @end_date::int
WHERE
    id = @campaign_id::uuid
RETURNING *;

-- name: UpdateCampaignTargeting :one
UPDATE campaigns_targeting
SET
    gender = COALESCE(sqlc.narg(gender)::varchar, NULL),
    age_from = COALESCE(sqlc.narg(age_from)::int, NULL), age_to = COALESCE(sqlc.narg(age_to)::int, NULL),
    location = COALESCE(sqlc.narg(location)::varchar, NULL)
WHERE
    campaign_id = @campaign_id::uuid
RETURNING *;

-- name: DeleteCampaignByID :exec
DELETE FROM campaigns
WHERE id = @campaign_id::uuid;

-- name: GetRelativeAd :one
SELECT * FROM campaigns 
JOIN campaigns_targeting ON campaigns.id = campaigns_targeting.campaign_id
JOIN ml_scores ON campaigns.advertiser_id = ml_scores.advertiser_id
WHERE 
    NOT EXISTS (
        SELECT 1 FROM impressions 
        WHERE impressions.campaign_id = campaigns.id 
          AND impressions.client_id = @client_id::uuid
    ) AND
    (
        SELECT COUNT(*) FROM impressions 
        WHERE impressions.campaign_id = campaigns.id
    ) < campaigns.impressions_limit
    AND
    (gender = @gender::varchar OR gender = 'ALL' OR gender IS NULL) AND
    (
        ((age_from IS NULL AND age_to >= @age::int) OR 
         (age_to IS NULL AND age_from <= @age::int) OR 
         (age_from IS NULL AND age_to IS NULL)) OR
        (age_from <= @age::int AND age_to >= @age::int)
    ) AND
    (location IS NULL OR location = @location::varchar) AND
    ml_scores.client_id = @client_id::uuid AND
    campaigns.start_date <= @cur_date::int AND 
    campaigns.end_date >= @cur_date::int
ORDER BY score DESC, cost_per_impression DESC
LIMIT 1;
