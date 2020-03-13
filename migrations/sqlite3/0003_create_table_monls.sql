-- +migration Up
CREATE TABLE IF NOT EXISTS monls (
	id          VARCHAR(50) PRIMARY KEY,
	url         VARCHAR(2000),
	title       VARCHAR(250),
	description TEXT,
	readme      TEXT,
	image_id    VARCHAR(50),
	created_at  TIMESTAMP,
	updated_at  TIMESTAMP,
	UNIQUE(url COLLATE NOCASE)
);

-- +migration Down
DROP TABLE IF EXISTS monls;
