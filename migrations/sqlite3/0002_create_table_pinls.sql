-- +migration Up
CREATE TABLE IF NOT EXISTS pinls (
	id VARCHAR(50) NOT NULL PRIMARY KEY,
	user_id VARCHAR(50) NOT NULL,
	url VARCHAR(255) NOT NULL,
	title VARCHAR(255) NOT NULL,
	description VARCHAR(255) NOT NULL,
	readme TEXT NOT NULL,
	image_id VARCHAR(50) NOT NULL,
	created_at TIMESTAMP NULL,
	updated_at TIMESTAMP NULL
);

-- +migration Down
DROP TABLE IF EXISTS pinls;
