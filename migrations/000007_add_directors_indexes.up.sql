CREATE INDEX IF NOT EXISTS directors_name_idx ON directors USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS directors_surname_idx ON directors USING GIN (to_tsvector('simple', surname));
CREATE INDEX IF NOT EXISTS directors_awards_idx ON directors USING GIN (awards);