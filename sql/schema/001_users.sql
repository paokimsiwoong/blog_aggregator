-- +goose Up
CREATE TABLE users(
    id uuid DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL,
    PRIMARY KEY (id)
);

-- @@@ 해답
-- CREATE TABLE users (
--     id UUID PRIMARY KEY,
--     created_at TIMESTAMP NOT NULL,
--     updated_at TIMESTAMP NOT NULL,
--     name TEXT NOT NULL UNIQUE
-- );

-- +goose Down
DROP TABLE users;
