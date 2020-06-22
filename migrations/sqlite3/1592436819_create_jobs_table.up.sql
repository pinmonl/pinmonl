CREATE TABLE IF NOT EXISTS jobs (
  id          VARCHAR(50) PRIMARY KEY,
  name        VARCHAR(100),
  describe    VARCHAR(250),
  target_id   VARCHAR(50),
  target_name VARCHAR(100),
  status      INTEGER,
  message     TEXT,
  created_at  TIMESTAMP,
  ended_at    TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ix_jobs_name ON jobs (name);
CREATE INDEX IF NOT EXISTS ix_jobs_target ON jobs (target_id, target_name);
CREATE INDEX IF NOT EXISTS ix_jobs_ended ON jobs (ended_at);
