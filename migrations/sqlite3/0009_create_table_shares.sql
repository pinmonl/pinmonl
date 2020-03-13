-- +migration Up
CREATE TABLE IF NOT EXISTS shares (
	id          VARCHAR(50) PRIMARY KEY,
	user_id     VARCHAR(50),
	name        VARCHAR(250),
	description TEXT,
	readme      TEXT,
	image_id    VARCHAR(50),
	created_at  TIMESTAMP,
	updated_at  TIMESTAMP
);

CREATE INDEX ix_share_user ON shares (user_id);
CREATE INDEX ix_share_name ON shares (name);

-- +migration Down
DROP TABLE IF EXISTS shares;
