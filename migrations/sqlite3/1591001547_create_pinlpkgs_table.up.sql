CREATE TABLE IF NOT EXISTS pinlpkgs (
  id      VARCHAR(50) PRIMARY KEY,
  pinl_id VARCHAR(50),
  pkg_id  VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS ix_pinlpkgs_pinl ON pinlpkgs (pinl_id);
CREATE INDEX IF NOT EXISTS ix_pinlpkgs_pkg ON pinlpkgs (pkg_id);
