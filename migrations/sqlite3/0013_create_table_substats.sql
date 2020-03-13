-- +migration Up
CREATE TABLE IF NOT EXISTS substats (
  id      VARCHAR(50) PRIMARY KEY,
  stat_id VARCHAR(50),
  kind    VARCHAR(50),
  labels  TEXT
);

CREATE INDEX ix_substat_stat ON substats (stat_id);
CREATE INDEX ix_substat_kind ON substats (kind);

-- +migration Down
DROP TABLE IF EXISTS substats;
