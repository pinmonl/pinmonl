-- +migration Up
CREATE TABLE IF NOT EXISTS pinmons (
	pinl_id VARCHAR(50) NOT NULL,
	monl_id VARCHAR(50) NOT NULL,
	user_id VARCHAR(50) NOT NULL,
	sort INTEGER NOT NULL
);

-- +migration Down
DROP TABLE IF EXISTS pinmons;
