-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- -- name: GetUser :one
-- SELECT * FROM users where name = $1 LIMIT 1;

-- -- name: DeleteUsers :exec
-- DELETE FROM users;

-- -- name: GetUsers :many
-- SELECT * FROM users;