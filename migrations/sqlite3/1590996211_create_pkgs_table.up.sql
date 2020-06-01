CREATE TABLE IF NOT EXISTS pkgs (
  id            VARCHAR(50) PRIMARY KEY,
  url           VARCHAR(2000),
  provider      VARCHAR(100),
  provider_host VARCHAR(100),
  provider_uri  VARCHAR(1000),
  created_at    TIMESTAMP,
  updated_at    TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS ix_pkgs_keys ON pkgs (provider, provider_host, provider_uri);
CREATE INDEX IF NOT EXISTS ix_pkgs_provider ON pkgs (provider);
CREATE INDEX IF NOT EXISTS ix_pkgs_provider_uri ON pkgs (provider_uri);
CREATE INDEX IF NOT EXISTS ix_pkgs_provider_host ON pkgs (provider_host);
