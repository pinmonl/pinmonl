-- +migration Up
CREATE TABLE IF NOT EXISTS images (
	id           VARCHAR(50) PRIMARY KEY,
	target_id    VARCHAR(50),
	target_name  VARCHAR(100),
	content_type VARCHAR(100),
	sort         INTEGER,
	filename     VARCHAR(250),
	content      BLOB,
	description  VARCHAR(250),
	size         INTEGER,
	created_at   TIMESTAMP,
	updated_at   TIMESTAMP
);

CREATE INDEX ix_image_target ON images (target_id, target_name);
CREATE INDEX ix_image_ctype ON images (content_type);

-- +migration Down
DROP TABLE IF EXISTS images;
