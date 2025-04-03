-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1, 
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)   
RETURNING *; 

-- @@@ 해답처럼 feeds 테이블도 추가로 join해서 feed name도 정보 받아오기
-- name: GetPostsForUser :many
SELECT posts.*, feeds.name AS feed_name FROM posts
-- WHERE feed_id IN (SELECT feed_id FROM feed_follows WHERE feed_follows.user_id = $1) -- deepseek : 간단하지만 (sub query)안에서 많은 데이터를 찾아오는 경우 성능 문제
INNER JOIN feed_follows
ON posts.feed_id = feed_follows.feed_id
INNER JOIN feeds
ON posts.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC NULLS LAST, posts.updated_at DESC
LIMIT $2;