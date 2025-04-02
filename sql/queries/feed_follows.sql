-- name: CreateFeedFollow :one
WITH returned AS (
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1, 
    $2,
    $3,
    $4,
    $5
)   
RETURNING *
)
SELECT returned.id, returned.created_at, returned.updated_at, returned.user_id, users.name AS user_name, returned.feed_id, feeds.name AS feed_name 
FROM returned
INNER JOIN users
ON returned.user_id = users.id
INNER JOIN feeds
ON returned.feed_id = feeds.id; 


-- @@@ :many를 해야 복수 결과 반환
-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, users.name AS user_name, feeds.name AS feed_name, feeds.url FROM feed_follows
INNER JOIN users
ON feed_follows.user_id = users.id
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;


-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = (SELECT id FROM users WHERE users.name = $1) AND feed_id = (SELECT id FROM feeds WHERE feeds.url = $2);