-- +migration Up
CREATE TABLE IF NOT EXISTS images (
	id VARCHAR(50) NOT NULL PRIMARY KEY,
	target_id VARCHAR(50) NOT NULL,
	target_name VARCHAR(100) NOT NULL,
	kind VARCHAR(100) NOT NULL,
	sort INTEGER,
	filename VARCHAR(255) NOT NULL,
	content BLOB NULL,
	description VARCHAR(255) NOT NULL,
	size INTEGER,
	created_at TIMESTAMP NULL,
	updated_at TIMESTAMP NULL
);

-- +migration Down
DROP TABLE IF EXISTS images;
