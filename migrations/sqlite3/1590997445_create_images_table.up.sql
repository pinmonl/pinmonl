CREATE TABLE IF NOT EXISTS images (
  id           VARCHAR(50) PRIMARY KEY,
  target_id    VARCHAR(50),
  target_name  VARCHAR(100),
  content      BLOB,
  description  VARCHAR(250),
  size         INTEGER,
  content_type VARCHAR(100),
  created_at   TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ix_images_target ON images (target_id, target_name);
CREATE INDEX IF NOT EXISTS ix_images_target ON images (content_type);
