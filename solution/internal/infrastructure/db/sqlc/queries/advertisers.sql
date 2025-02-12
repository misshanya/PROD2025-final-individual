-- name: CreateAdvertiser :exec
INSERT INTO advertisers (
    id, name
) VALUES (
    @id::uuid, @name::varchar
);

-- name: GetAdvertiserByID :one
SELECT * FROM advertisers
WHERE id = @id::uuid;
