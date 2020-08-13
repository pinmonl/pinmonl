CREATE TABLE IF NOT EXISTS shares (
  id          VARCHAR(50) PRIMARY KEY,
  user_id     VARCHAR(50),
  slug        VARCHAR(250),
  name        VARCHAR(250),
  description TEXT,
  image_id    VARCHAR(50),
  status      INTEGER,
  created_at  TIMESTAMP,
  updated_at  TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ix_shares_user ON shares (user_id);
CREATE INDEX IF NOT EXISTS ix_shares_slug ON shares (slug);
