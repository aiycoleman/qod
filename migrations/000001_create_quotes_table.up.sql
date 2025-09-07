--Filename: migrations/00001_create_quotes_table.up.sql
CREATE TABLE IF NOT EXISTS quotes (
    id bigserial PRIMARY KEY,
    content text NOT NULL,
    author text NOT NULL,
    version integer NOT NULL DEFAULT 1,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);