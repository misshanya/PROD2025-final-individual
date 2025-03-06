-- name: CreateAdvertiser :exec
INSERT INTO advertisers (
    id, name
) VALUES (
    @id::uuid, @name::varchar
);

-- name: GetAdvertiserByID :one
SELECT * FROM advertisers
WHERE id = @id::uuid;

-- name: UpdateAdvertiser :exec
UPDATE advertisers
SET name = @name::varchar
WHERE id = @id::uuid;

-- name: CreateMLScore :exec
INSERT INTO ml_scores (
    client_id, advertiser_id, score
) VALUES (
    @client_id::uuid, @advertiser_id::uuid, @score::int
);

-- name: GetMLScoreByIDs :one
SELECT * FROM ml_scores
WHERE
client_id = @client_id::uuid AND
advertiser_id = @advertiser_id::uuid;

-- name: UpdateMLScore :exec
UPDATE ml_scores
SET score = @score::int
WHERE
client_id = @client_id::uuid AND
advertiser_id = @advertiser_id::uuid;