CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS followers (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    follower_id bigint NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);