-- +migration Up
CREATE TABLE IF NOT EXISTS sharetags (
	share_id  VARCHAR(50),
	tag_id    VARCHAR(50),
	kind      INTEGER,
	parent_id VARCHAR(50),
	sort      INTEGER,
	level     INTEGER
);

CREATE INDEX ix_sharetag_share ON sharetags (share_id);
CREATE INDEX ix_sharetag_kind ON sharetags (kind);
CREATE INDEX ix_sharetag_tag ON sharetags (tag_id);
CREATE INDEX ix_sharetag_parent ON sharetags (parent_id);
CREATE INDEX ix_sharetag_level ON sharetags (level);

-- +migration Down
DROP TABLE IF EXISTS sharetags;
