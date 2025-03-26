-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1, -- $1,2,3,4 는 go 코드에서 sql query의 value들을 채우기 위해 입력하는 변수들을 위한 place holder
    $2,
    $3,
    $4
)
RETURNING *; -- returning 으로 생성한 유저를 바로 반환하고 있음 (위에 :one으로 생성한 유저 하나만 반환하도록 함)

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1; -- name이 UNIQUE이므로 LIMIT 1 필요 없음


-- name: ResetUsers :exec
TRUNCATE TABLE users; -- :exec : The generated method will return the error from ExecContext.
-- TRUNCATE : https://knitter-amelie.tistory.com/21
-- @@@ 해답은 DELETE FROM users; 사용