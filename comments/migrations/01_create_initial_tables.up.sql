CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS POSTGIS;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS comments CASCADE;

create table comments
(
    comment_id UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    hotel_id   UUID   NOT NULL,
    user_id    UUID   NOT NULL,
    message    citext NOT NULL CHECK ( message <> '' ),
    photos     text[],
    rating     float                    DEFAULT 0 CHECK (rating >= 0 AND rating <= 10),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);