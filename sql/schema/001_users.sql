-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL
);
-- @@@ internal/database 안의 go파일들을 보면 uuid를 구글의 uuid패키지를 이용해서 생성해서 입력하고 있으므로 해답 방식으로 변경
-- CREATE TABLE users(
--     id uuid DEFAULT gen_random_uuid(), -- 이방식은 postgreSQL 기능으로 자동생성 ==> insert 할떄 id 필드 명시 안해도 됨
--     created_at TIMESTAMP NOT NULL,
--     updated_at TIMESTAMP NOT NULL,
--     name TEXT UNIQUE NOT NULL,
--     PRIMARY KEY (id)
-- );

-- @@@ 해답
-- CREATE TABLE users (
--     id UUID PRIMARY KEY,
--     created_at TIMESTAMP NOT NULL,
--     updated_at TIMESTAMP NOT NULL,
--     name TEXT NOT NULL UNIQUE
-- );

-- +goose Down
DROP TABLE users;
