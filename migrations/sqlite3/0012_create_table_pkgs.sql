-- +migration Up
CREATE TABLE IF NOT EXISTS pkgs (
	id          VARCHAR(50) PRIMARY KEY,
	monl_id     VARCHAR(50),
	url         VARCHAR(2000),
	vendor      VARCHAR(100),
	vendor_uri  VARCHAR(1000),
	title       VARCHAR(250),
	description TEXT,
	readme      TEXT,
	image_id    VARCHAR(50),
	labels      TEXT,
	created_at  TIMESTAMP,
	updated_at  TIMESTAMP
);

CREATE INDEX ix_pkg_monl ON pkgs (monl_id);
CREATE INDEX ix_pkg_vendor ON pkgs (vendor);
CREATE INDEX ix_pkg_vendor_uri ON pkgs (vendor, vendor_uri);

-- +migration Down
DROP TABLE IF EXISTS pkgs;
