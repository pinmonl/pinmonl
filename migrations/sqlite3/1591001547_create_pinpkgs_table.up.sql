CREATE TABLE IF NOT EXISTS pinpkgs (
  id      VARCHAR(50) PRIMARY KEY,
  pinl_id VARCHAR(50),
  pkg_id  VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS ix_pinpkgs_pinl ON pinpkgs (pinl_id);
CREATE INDEX IF NOT EXISTS ix_pinpkgs_pkg ON pinpkgs (pkg_id);
