CREATE TABLE IF NOT EXISTS sharetags (
  id           VARCHAR(50) PRIMARY KEY,
  share_id     VARCHAR(50),
  tag_id       VARCHAR(50),
  kind         INTEGER,
  parent_id    VARCHAR(50),
  level        INTEGER,
  status       INTEGER,
  has_children BOOLEAN
);

CREATE INDEX IF NOT EXISTS ix_sharetags_share ON sharetags (share_id);
CREATE INDEX IF NOT EXISTS ix_sharetags_tag ON sharetags (tag_id);
