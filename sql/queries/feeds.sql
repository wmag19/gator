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

-- name: GetFeedsAndUser :many
SELECT * FROM feeds JOIN users ON feeds.user_id = users.id;
-- -- name: GetFeeds :many
-- SELECT * FROM feeds JOIN users ON feeds.user_id = users.id;
-- -- name: GetFeedsFromUser :many
-- SELECT * FROM feeds WHERE user_id = $1;
-- -- name: GetUserFromFeed :one
-- SELECT * FROM feeds JOIN users ON feeds.user_id = users.id WHERE feeds.id = $1;

-- name: GetUserNameFromFeedID

-- -- name: DeleteUsers :exec
-- DELETE FROM users;

-- -- name: GetUsers :many
-- SELECT * FROM users;