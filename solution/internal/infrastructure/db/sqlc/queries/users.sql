-- name: CreateUser :one
INSERT INTO users (
    id, login, age, location, gender
) VALUES (
    @id::uuid, @login::varchar,
    @age::integer, @location::varchar, @gender::varchar
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = @id::uuid;

-- name: UpdateUser :exec
UPDATE users
SET
login = @login::varchar,
age = @age::int,
location = @location::varchar,
gender = @gender::varchar
WHERE id = @id::uuid;