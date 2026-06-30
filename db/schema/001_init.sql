CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    director TEXT NOT NULL,
    release_year INT NOT NULL,
    genre TEXT NOT NULL,
    available_copies INT NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS loans (
    id SERIAL PRIMARY KEY,
    movie_id INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    loan_date TIMESTAMP NOT NULL DEFAULT NOW(),
    return_date TIMESTAMP,
    price DECIMAL(10, 2) NOT NULL,
    returned BOOLEAN NOT NULL DEFAULT FALSE
);

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'loans'
        AND column_name = 'return_date'
        AND is_nullable = 'NO'
    ) THEN
        ALTER TABLE loans ALTER COLUMN return_date DROP NOT NULL;
    END IF;
END$$;