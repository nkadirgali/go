CREATE TABLE IF NOT EXISTS actors (
    -- id column is a 64-bit auto-incrementing integer & primary key (defines the row)
    id bigserial PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    version integer NOT NULL DEFAULT 1
);

