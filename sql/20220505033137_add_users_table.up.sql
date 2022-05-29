CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  email TEXT NOT NULL,
  verified BOOLEAN NOT NULL,
  created_at TIMESTAMP DEFAULT current_timestamp
);