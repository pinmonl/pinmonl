-- +migration Up
CREATE TABLE IF NOT EXISTS pkgs (
	id            VARCHAR(50) PRIMARY KEY,
	url           VARCHAR(2000),
	provider      VARCHAR(100),
	provider_host VARCHAR(100),
	provider_uri  VARCHAR(1000),
	title         VARCHAR(250),
	description   TEXT,
	readme        TEXT,
	image_id      VARCHAR(50),
	labels        TEXT,
	created_at    TIMESTAMP,
	updated_at    TIMESTAMP,
	UNIQUE(provider, provider_host, provider_uri)
);

CREATE INDEX ix_pkg_provider ON pkgs (provider);
CREATE INDEX ix_pkg_provider_host ON pkgs (provider, provider_host);
CREATE INDEX ix_pkg_provider_uri ON pkgs (provider, provider_uri);

-- +migration Down
DROP TABLE IF EXISTS pkgs;
