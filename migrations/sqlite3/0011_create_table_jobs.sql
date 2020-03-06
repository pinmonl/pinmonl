-- +migration Up
CREATE TABLE IF NOT EXISTS jobs (
	id VARCHAR(50) NOT NULL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	target_id VARCHAR(50) NOT NULL,
	status VARCHAR(100) NOT NULL,
	error VARCHAR(255) NOT NULL,
	scheduled_at TIMESTAMP NULL,
	started_at TIMESTAMP NULL,
	created_at TIMESTAMP NULL
);

-- +migration Down
DROP TABLE IF EXISTS jobs;
