CREATE TABLE IF NOT EXISTS users (
  id         VARCHAR(50) PRIMARY KEY,
  login      VARCHAR(250),
  password   VARCHAR(250),
  name       VARCHAR(250),
  image_id   VARCHAR(50),
  hash       VARCHAR(500),
  role       INTEGER,
  status     INTEGER,
  last_seen  TIMESTAMP,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ix_users_login ON users (login COLLATE NOCASE);
CREATE INDEX IF NOT EXISTS ix_users_role ON users (role);
CREATE INDEX IF NOT EXISTS ix_users_hash ON users (hash);
CREATE INDEX IF NOT EXISTS ix_users_status ON users (status);
