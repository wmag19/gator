-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at = $1, updated_at = $2 WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT 
    "id", 
    "name", 
    "url", 
    "last_fetched_at" 
FROM 
    feeds 
ORDER BY 
    "last_fetched_at" ASC NULLS FIRST 
LIMIT 1;
