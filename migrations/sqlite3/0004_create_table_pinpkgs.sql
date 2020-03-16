-- +migration Up
CREATE TABLE IF NOT EXISTS pinpkgs (
	pinl_id VARCHAR(50),
	pkg_id  VARCHAR(50),
	sort    INTEGER,
	UNIQUE(pinl_id, pkg_id)
);

CREATE INDEX ix_pinpkg_pinl ON pinpkgs (pinl_id);
CREATE INDEX ix_pinpkg_pkg ON pinpkgs (pkg_id);

-- +migration Down
DROP TABLE IF EXISTS pinpkgs;
