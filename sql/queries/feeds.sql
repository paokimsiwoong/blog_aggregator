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
RETURNING *; -- returning 으로 생성한 유저를 바로 반환하고 있음 (위에 :one으로 생성한 유저 하나만 반환하도록 함)

-- name: GetFeed :one
SELECT * FROM feeds
WHERE name = $1; -- name이 UNIQUE이므로 LIMIT 1 필요 없음


-- name: GetFeeds :many
SELECT * FROM feeds;