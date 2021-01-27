DROP TABLE IF EXISTS users CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;
-- CREATE EXTENSION IF NOT EXISTS postgis;
-- CREATE EXTENSION IF NOT EXISTS postgis_topology;


CREATE TYPE role AS ENUM ('admin', 'user', 'owner');

CREATE TABLE users
(
    user_id    UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    first_name VARCHAR(60)        NOT NULL CHECK ( first_name <> '' ),
    last_name  VARCHAR(60)        NOT NULL CHECK ( last_name <> '' ),
    email      VARCHAR(64) UNIQUE NOT NULL CHECK ( email <> '' ),
    password   TEXT               NOT NULL CHECK ( octet_length(password) <> 0 ),
    role       role                     DEFAULT 'user',
    avatar     TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION user_updated() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_updated_at_trigger
    BEFORE INSERT OR UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE user_updated();