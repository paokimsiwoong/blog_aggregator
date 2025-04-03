-- returning 으로 생성한 유저를 바로 반환하고 있음 (위에 :one으로 생성한 유저 하나만 반환하도록 함)
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

-- name이 UNIQUE이므로 LIMIT 1 필요 없음
-- name: GetFeed :one
SELECT * FROM feeds 
WHERE name = $1; 

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;


-- @@@ 해답은 -- name: GetFeeds :many <br> SELECT * FROM feeds; 로 간단하게 한 뒤
-- @@@ users.sql에서 -- name: GetUserById :one <br> SELECT * FROM users WHERE id = $1; 을 추가하고 
-- @@@ commands.go 에서 feeds에 저장된 user_id를 GetUserById 함수에 입력해 User 구조체를 불러내어 사용
-- name: GetFeeds :many
SELECT feeds.id, feeds.name, feeds.created_at, feeds.updated_at, feeds.url, users.name AS user_name
FROM feeds
INNER JOIN users
ON feeds.user_id = users.id
ORDER BY feeds.updated_at;


-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $1, last_fetched_at = $1 -- @@@ 해답은 NOW() 함수로 sql 안에서 해결
WHERE id = $2;


-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
