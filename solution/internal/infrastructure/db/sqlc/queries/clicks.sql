-- name: CreateClick :one
INSERT INTO clicks (
    campaign_id, client_id
) VALUES (
    @campaign_id::uuid, @client_id::uuid
)
RETURNING *;

-- name: IsClicked :one
SELECT 1 FROM clicks
WHERE
    campaign_id = @campaign_id::uuid AND
    client_id = @client_id::uuid;