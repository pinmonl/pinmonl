CREATE TABLE IF NOT EXISTS sharepins (
  id       VARCHAR(50) PRIMARY KEY,
  share_id VARCHAR(50),
  pinl_id  VARCHAR(50),
  status   INTEGER
);

CREATE INDEX IF NOT EXISTS ix_sharepins_share ON sharepins (share_id);
CREATE INDEX IF NOT EXISTS ix_sharepins_pinl ON sharepins (pinl_id);
