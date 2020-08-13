CREATE TABLE IF NOT EXISTS tags (
  id           VARCHAR(50) PRIMARY KEY,
  name         VARCHAR(250),
  user_id      VARCHAR(50),
  parent_id    VARCHAR(50),
  level        INTEGER,
  color        VARCHAR(20),
  bg_color     VARCHAR(20),
  has_children BOOLEAN,
  created_at   TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ix_tags_user ON tags (user_id);
CREATE INDEX IF NOT EXISTS ix_tags_level ON tags (level);
CREATE INDEX IF NOT EXISTS ix_tags_parent ON tags (parent_id);
