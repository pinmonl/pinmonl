CREATE TABLE IF NOT EXISTS monpkgs (
  id      VARCHAR(50) PRIMARY KEY,
  monl_id VARCHAR(50),
  pkg_id  VARCHAR(50),
  kind    INTEGER
);

CREATE UNIQUE INDEX IF NOT EXISTS ix_monpkgs_keys ON monpkgs (monl_id, pkg_id);
