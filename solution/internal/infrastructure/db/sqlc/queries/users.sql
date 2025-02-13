-- name: CreateUser :one
INSERT INTO users (
    id, login, age, location, gender
) VALUES (
    @id::uuid, @login::varchar(50),
    @age::integer, @location::varchar, @gender::varchar(7)
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
gender = @gender::varchar(7)
WHERE id = @id::uuid;