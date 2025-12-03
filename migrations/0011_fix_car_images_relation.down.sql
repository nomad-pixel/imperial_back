ALTER TABLE cars ADD COLUMN image_url TEXT;

DROP INDEX IF EXISTS idx_car_images_car_id;

ALTER TABLE car_images DROP COLUMN car_id;
