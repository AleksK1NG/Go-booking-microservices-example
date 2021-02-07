-- ALTER SYSTEM SET shared_buffers = '128MB';
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS POSTGIS;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS hotels CASCADE;

create table hotels
(
    hotel_id       UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    name           citext NOT NULL CHECK ( name <> '' ),
    email          citext NOT NULL CHECK ( email <> '' ),
    location       citext NOT NULL CHECK ( location <> '' ),
    country        citext NOT NULL CHECK ( country <> '' ),
    city           citext NOT NULL CHECK ( city <> '' ),
    coordinates    geometry(POINT, 4326),
    description    text,
    image          text,
    photos         text[],
    rating         float                    DEFAULT 0 CHECK (rating >= 0 AND rating <= 10),
    comments_count int                      DEFAULT 0 CHECK (comments_count >= 0),
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX if not exists hotels_gist_idx ON hotels USING gist (coordinates);

CREATE INDEX hotels_name_trgm_idx ON hotels
    USING gist (name);

CREATE INDEX hotels_location_trgm_idx ON hotels
    USING gist (location);

CREATE OR REPLACE FUNCTION hotel_updated() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER hotel_updated_at_trigger
    BEFORE INSERT OR UPDATE
    ON hotels
    FOR EACH ROW
EXECUTE PROCEDURE hotel_updated();