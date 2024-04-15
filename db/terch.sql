CREATE TABLE IF NOT EXISTS docs (
    id bigserial PRIMARY KEY,
    name text NOT NULL, 
    docvec DOUBLE PRECISION[] NOT NULL 
); 