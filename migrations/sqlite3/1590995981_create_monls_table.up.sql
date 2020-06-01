CREATE TABLE IF NOT EXISTS monls (
  id         VARCHAR(50) PRIMARY KEY,
  url        VARCHAR(2000),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS ix_monls_url ON monls (url COLLATE NOCASE);
