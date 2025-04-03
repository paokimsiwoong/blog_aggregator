-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP;
-- @@@ 해답처럼 별개의 sql 파일로 db 변경사항이 차곡차곡 쌓이도록 변경


-- +goose Down
ALTER TABLE feeds
DROP COLUMN last_fetched_at;