-- +migration Up
CREATE TABLE IF NOT EXISTS monpkgs (
  monl_id VARCHAR(50),
  pkg_id  VARCHAR(50),
  tie     INTEGER,
  UNIQUE(monl_id, pkg_id)
);

CREATE INDEX ix_monpkg_monl ON monpkgs (monl_id);
CREATE INDEX ix_monpkg_pkg ON monpkgs (pkg_id);

-- +migration Down
DROP TABLE IF EXISTS monpkgs;
