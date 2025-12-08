-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8 );

-- name: GetPostsForUser :many
SELECT 
    p."id", 
    p."title", 
    p."url", 
    p."description", 
    p."published_at" 
FROM 
    posts p
JOIN 
    feeds f ON p."feed_id" = f."id"
WHERE 
    f."user_id" = $1
ORDER BY 
    p."published_at" DESC
LIMIT $2;

-- -- name: GetUser :one
-- SELECT * FROM users where name = $1 LIMIT 1;

-- -- name: DeleteUsers :exec
-- DELETE FROM users;

-- -- name: GetUsers :many
-- SELECT * FROM users;