-- name: CreateFeedFollows :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id) 
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollows :many 
SELECT * from feed_follows where user_id = $1;

-- name: DeleteFeedFollow :exec 
DELETE from feed_follows where id = $1 and user_id = $2;