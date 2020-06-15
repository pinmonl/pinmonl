CREATE TABLE IF NOT EXISTS pinls (
  id          VARCHAR(50) PRIMARY KEY,
  user_id     VARCHAR(50),
  monl_id     VARCHAR(50),
  url         VARCHAR(2000),
  title       VARCHAR(250),
  description TEXT,
  image_id    VARCHAR(50),
  status      INTEGER,
  created_at  TIMESTAMP,
  updated_at  TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ix_pinls_user ON pinls (user_id);
CREATE INDEX IF NOT EXISTS ix_pinls_monl ON pinls (monl_id);
