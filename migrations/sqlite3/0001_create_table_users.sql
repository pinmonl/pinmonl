-- +migration Up
CREATE TABLE IF NOT EXISTS users (
	id         VARCHAR(50) PRIMARY KEY,
	login      VARCHAR(250),
	password   VARCHAR(250),
	name       VARCHAR(250),
	image_id   VARCHAR(50),
	role       INTEGER,
	hash       VARCHAR(500),
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	last_log   TIMESTAMP,
	UNIQUE(login COLLATE NOCASE)
);

CREATE INDEX ix_user_role ON users (role);
CREATE INDEX ix_user_hash ON users (hash);

-- +migration Down
DROP TABLE IF EXISTS users;
