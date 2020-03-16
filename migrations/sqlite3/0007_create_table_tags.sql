-- +migration Up
CREATE TABLE IF NOT EXISTS tags (
	id         VARCHAR(50) PRIMARY KEY,
	name       VARCHAR(250),
	user_id    VARCHAR(50),
	parent_id  VARCHAR(50),
	sort       INTEGER,
	level      INTEGER,
	color      VARCHAR(50),
	bgcolor    VARCHAR(50),
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);

CREATE INDEX ix_tag_user ON tags (user_id);
CREATE INDEX ix_tag_user_name ON tags (user_id, name);
CREATE INDEX ix_tag_parent ON tags (parent_id);
CREATE INDEX ix_tag_level ON tags (level);

-- +migration Down
DROP TABLE IF EXISTS tags;
