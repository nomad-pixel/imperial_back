DELETE FROM car_images;

ALTER TABLE car_images ADD COLUMN car_id INT NOT NULL REFERENCES cars(id) ON DELETE CASCADE;

CREATE INDEX idx_car_images_car_id ON car_images(car_id);

ALTER TABLE cars DROP COLUMN IF EXISTS image_url;
