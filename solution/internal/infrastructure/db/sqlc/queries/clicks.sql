-- name: CreateClick :one
INSERT INTO clicks (
    campaign_id, client_id
) VALUES (
    @campaign_id::uuid, @client_id::uuid
)
RETURNING *;
