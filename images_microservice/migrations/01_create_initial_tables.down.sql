DROP TABLE IF EXISTS images CASCADE;
DROP TRIGGER IF EXISTS image_updated_at_trigger ON images;
DROP FUNCTION IF EXISTS image_updated;