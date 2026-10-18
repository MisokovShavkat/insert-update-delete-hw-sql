CREATE TABLE IF NOT EXISTS books(
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    author TEXT NOT NULL,
    price INT NOT NULL,
    pages INT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);