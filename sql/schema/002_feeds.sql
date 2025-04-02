-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL, -- 동일 url 중복을 막아서 같은 사이트에서 같은 내용 rss feed 안받게 막기
    user_id UUID NOT NULL,
    CONSTRAINT fk_users
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE 
    -- ON DELETE CASCADE, ON UPDATE CASCADE 등은 user_id가 있는 원본 테이블에서 해당 record가 변동되었을때(delete나 update) 
    -- 그 해당 record의 user_id를 foreign key로 가지고 있는 feeds 테이블의 모든 record에도 변동 결과를 적용한다.
    -- @@@ 해답은 user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE 
);

-- +goose Down
DROP TABLE feeds;