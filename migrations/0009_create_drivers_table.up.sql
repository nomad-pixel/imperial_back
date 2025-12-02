Create TAble IF NOT EXISTS drivers (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    about VARCHAR(1000) NOT NULL,
    photo_url VARCHAR(500) NOT NULL,
    experience_years VARCHAR(300) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()

);