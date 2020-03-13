-- +migration Up
CREATE TABLE IF NOT EXISTS taggables (
	tag_id      VARCHAR(50),
	target_id   VARCHAR(50),
	target_name VARCHAR(100),
	sort        INTEGER
);

CREATE INDEX ix_taggable_tag ON taggables (tag_id);
CREATE INDEX ix_taggable_target ON taggables (target_id, target_name);

-- +migration Down
DROP TABLE IF EXISTS taggables;
