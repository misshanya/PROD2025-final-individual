-- name: CreateImpression :one
INSERT INTO impressions (
    campaign_id, client_id
) VALUES (
    @campaign_id::uuid, @client_id::uuid
)
RETURNING *;
