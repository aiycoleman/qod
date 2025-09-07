-- Filename: migrations/000002_add_quotes_likes_column.up.sql
ALTER TABLE quotes
ADD COLUMN likes integer NOT NULL DEFAULT 0;

