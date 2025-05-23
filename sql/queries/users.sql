-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1, -- $1,2,3,4 는 go 코드에서 sql query의 value들을 채우기 위해 입력하는 변수들을 위한 place holder (@@@ postgreSQL만 $1, $2, ... 사용. MySQL은 ?, Oracle은 :val1, :val2 사용)
    $2,
    $3,
    $4
)
RETURNING *; -- returning 으로 생성한 유저를 바로 반환하고 있음 (위에 :one으로 생성한 유저 하나만 반환하도록 함)

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1; -- name이 UNIQUE이므로 LIMIT 1 필요 없음


-- name: GetUsers :many
SELECT * FROM users;

-- name: ResetUsers :exec
TRUNCATE TABLE users CASCADE; -- :exec : The generated method will return the error from ExecContext.
-- TRUNCATE : https://knitter-amelie.tistory.com/21
-- @@@ 해답은 DELETE FROM users; 사용 + CASCADE 명령어 없어도 문제 없음 
-- @@@ (DELETE는 한줄한줄 지우면서 기본 설정으로 CASCADE 반영 => CASCADE 명시 안해도 됨)
-- @@@ (TRUNCATE는 기본 설정으로 CASCADE 반영하지 않음 ==> CASCADE 명시 필수)