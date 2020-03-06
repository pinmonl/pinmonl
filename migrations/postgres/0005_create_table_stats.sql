-- +migration Up
CREATE TABLE IF NOT EXISTS stats (
	id VARCHAR(50) NOT NULL PRIMARY KEY,
	pkg_id VARCHAR(50) NOT NULL,
	recorded_at TIMESTAMP NULL,
	kind VARCHAR(100) NOT NULL,
	value VARCHAR(255) NOT NULL,
	is_latest BOOLEAN,
	manifest TEXT NOT NULL
);

-- +migration Down
DROP TABLE IF EXISTS stats;
