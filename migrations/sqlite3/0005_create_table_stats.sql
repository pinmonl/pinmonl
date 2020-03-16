-- +migration Up
CREATE TABLE IF NOT EXISTS stats (
	id          VARCHAR(50) PRIMARY KEY,
	pkg_id      VARCHAR(50),
	recorded_at TIMESTAMP,
	kind        VARCHAR(100),
	value       VARCHAR(250),
	is_latest   BOOLEAN,
	labels      TEXT
);

CREATE INDEX ix_stat_pkg ON stats (pkg_id);
CREATE INDEX ix_stat_latest ON stats (is_latest);
CREATE INDEX ix_stat_kind ON stats (kind);
CREATE INDEX ix_stat_recorded ON stats (recorded_at);
CREATE INDEX ix_stat_pkg_value ON stats (pkg_id, value);

-- +migration Down
DROP TABLE IF EXISTS stats;
