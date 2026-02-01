ALTER TABLE jobs RENAME published_at TO created_at;
ALTER TABLE jobs
    ADD COLUMN IF NOT EXISTS published BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS published_at TIMESTAMP DEFAULT null;
