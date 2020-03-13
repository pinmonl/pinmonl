-- +migration Up
CREATE TABLE IF NOT EXISTS pinmons (
	user_id VARCHAR(50),
	pinl_id VARCHAR(50),
	pkg_id  VARCHAR(50),
	sort    INTEGER
);

CREATE INDEX ix_pinmon_user ON pinmons (user_id);
CREATE INDEX ix_pinmon_pinl ON pinmons (pinl_id);
CREATE INDEX ix_pinmon_pkg ON pinmons (pkg_id);

-- +migration Down
DROP TABLE IF EXISTS pinmons;
