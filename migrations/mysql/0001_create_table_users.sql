-- +migration Up
CREATE TABLE IF NOT EXISTS users (
	id VARCHAR(50) NOT NULL,
	login VARCHAR(255) NOT NULL,
	password VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	image_id VARCHAR(50) NOT NULL,
	created_at TIMESTAMP NULL,
	updated_at TIMESTAMP NULL
);

-- +migration Down
DROP TABLE IF EXISTS users;
