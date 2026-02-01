ALTER TABLE jobs
    DROP COLUMN IF EXISTS published,
    DROP COLUMN IF EXISTS published_at;
ALTER TABLE jobs RENAME created_at TO published_at;
