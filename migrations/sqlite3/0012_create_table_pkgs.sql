-- +migration Up
CREATE TABLE IF NOT EXISTS pkgs (
	id          VARCHAR(50) PRIMARY KEY,
	url         VARCHAR(2000),
	vendor      VARCHAR(100),
	vendor_uri  VARCHAR(1000),
	title       VARCHAR(250),
	description TEXT,
	readme      TEXT,
	image_id    VARCHAR(50),
	labels      TEXT,
	created_at  TIMESTAMP,
	updated_at  TIMESTAMP,
	UNIQUE(vendor, vendor_uri)
);

CREATE INDEX ix_pkg_vendor ON pkgs (vendor);

-- +migration Down
DROP TABLE IF EXISTS pkgs;
