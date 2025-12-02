CREATE TABLE IF NOT EXISTS leads (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    phone varchar(32),
    period tstzrange NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);