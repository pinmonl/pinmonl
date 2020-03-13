-- +migration Up
CREATE TABLE IF NOT EXISTS pinls (
	id          VARCHAR(50) PRIMARY KEY,
	user_id     VARCHAR(50),
	url         VARCHAR(2000),
	title       VARCHAR(250),
	description TEXT,
	readme      TEXT,
	image_id    VARCHAR(50),
	created_at  TIMESTAMP,
	updated_at  TIMESTAMP
);

CREATE INDEX ix_pinl_user ON pinls (user_id);

-- +migration Down
DROP TABLE IF EXISTS pinls;
