CREATE TABLE IF NOT EXISTS stats (
  id           VARCHAR(50) PRIMARY KEY,
  pkg_id       VARCHAR(50),
  parent_id    VARCHAR(50),
  recorded_at  TIMESTAMP,
  kind         INTEGER,
  value        VARCHAR(250),
  value_type   INTEGER,
  checksum     VARCHAR(500),
  weight       INTEGER,
  is_latest    BOOLEAN,
  has_children BOOLEAN
);

CREATE INDEX IF NOT EXISTS ix_stats_pkg ON stats (pkg_id);
CREATE INDEX IF NOT EXISTS ix_stats_latest ON stats (is_latest);
CREATE INDEX IF NOT EXISTS ix_stats_parent ON stats (parent_id);
CREATE INDEX IF NOT EXISTS ix_stats_kind ON stats (kind);
CREATE INDEX IF NOT EXISTS ix_stats_value_type ON stats (value_type);
