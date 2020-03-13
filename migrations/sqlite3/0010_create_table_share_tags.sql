-- +migration Up
CREATE TABLE IF NOT EXISTS share_tags (
	share_id  VARCHAR(50),
	tag_id    VARCHAR(50),
	kind      INTEGER,
	parent_id VARCHAR(50),
	sort      INTEGER,
	level     INTEGER
);

CREATE INDEX ix_share_tag_share ON share_tags (share_id);
CREATE INDEX ix_share_tag_kind ON share_tags (kind);
CREATE INDEX ix_share_tag_tag ON share_tags (tag_id);
CREATE INDEX ix_share_tag_parent ON share_tags (parent_id);
CREATE INDEX ix_share_tag_level ON share_tags (level);

-- +migration Down
DROP TABLE IF EXISTS share_tags;
