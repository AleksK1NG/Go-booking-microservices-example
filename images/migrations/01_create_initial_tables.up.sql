DROP TABLE IF EXISTS images CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE images
(
    image_id    UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    image_url   VARCHAR(250) NOT NULL CHECK ( image_url <> '' ),
    is_uploaded bool                     DEFAULT false,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION image_updated() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER image_updated_at_trigger
    BEFORE INSERT OR UPDATE
    ON images
    FOR EACH ROW
EXECUTE PROCEDURE image_updated();