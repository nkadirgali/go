CREATE TABLE IF NOT EXISTS directors (
    id bigserial PRIMARY KEY,
    name text not null,
    surname text not null,
    awards text[] not NULL
);

