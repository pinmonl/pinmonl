-- +migration Up
CREATE TABLE IF NOT EXISTS tags (
	id VARCHAR(50) NOT NULL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	user_id VARCHAR(50) NOT NULL,
	parent_id VARCHAR(50) NOT NULL,
	sort INTEGER,
	created_at TIMESTAMP NULL,
	updated_at TIMESTAMP NULL
);

-- +migration Down
DROP TABLE IF EXISTS tags;
