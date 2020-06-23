CREATE TABLE IF NOT EXISTS monls (
  id         VARCHAR(50) PRIMARY KEY,
  url        VARCHAR(2000),
  fetched_at TIMESTAMP,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS ix_monls_url ON monls (url COLLATE NOCASE);
CREATE UNIQUE INDEX IF NOT EXISTS ix_monls_fetched_at ON monls (fetched_at);
