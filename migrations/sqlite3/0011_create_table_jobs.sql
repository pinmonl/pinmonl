-- +migration Up
CREATE TABLE IF NOT EXISTS jobs (
	id           VARCHAR(50) PRIMARY KEY,
	name         INTEGER,
	target_id    VARCHAR(50),
	status       INTEGER,
	error        TEXT,
	scheduled_at TIMESTAMP,
	started_at   TIMESTAMP,
	created_at   TIMESTAMP
);

CREATE INDEX ix_job_name ON jobs (name);
CREATE INDEX ix_job_status ON jobs (status);
CREATE INDEX ix_job_schedule ON jobs (scheduled_at);

-- +migration Down
DROP TABLE IF EXISTS jobs;
