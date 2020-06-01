CREATE TABLE IF NOT EXISTS taggables (
  id            VARCHAR(50) PRIMARY KEY,
  tag_id        VARCHAR(50),
  taggable_id   VARCHAR(50),
  taggable_name VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS ix_taggables_tag ON taggables (tag_id);
CREATE INDEX IF NOT EXISTS ix_taggables_target ON taggables (taggable_id, taggable_name);
