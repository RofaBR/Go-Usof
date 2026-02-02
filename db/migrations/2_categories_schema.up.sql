CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT
);

CREATE INDEX idx_categories_slug ON categories(slug);