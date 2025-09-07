--Filename: migrations/000002_add_quotes_likes_column.down.sql
ALTER TABLE quotes
DROP COLUMN IF EXISTS likes;