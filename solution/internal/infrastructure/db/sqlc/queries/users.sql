-- name: CreateUser :exec
INSERT INTO users (
    id, login, age, location, gender
) VALUES (
    @id::uuid, @login::varchar(50),
    @age::integer, @location::varchar, @gender::varchar(7)
)
RETURNING *;