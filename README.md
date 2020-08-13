# Pinmonl

Pinmonl is a bookmark manager with hierarchical tags which aims to be fast, shareable and monitor to releases.

My primary goal is to manage and monitor release of the git repositories from different providers (GitHub, GitLab, etc.) and their corresponding release in package managers (NPM, Docker, etc.). With sharing feature, it would be a better alternative for awesome list hosting.

Beside of Pinmonl, an Exchange server is developed to handle crawlers, notify new release to end-client and host sharing. The communication to Exchange server is an opt-in feature that is on by default. Release notification will send to the end-client whose with active connection to Exchange server only.

## Features

- Hierarchical tags
- Keyboard bindings
- Support SQLite and Postgres
- Custom thumbnail
- Show releases and statistical information if available
- Fill bookmark information by meta tags
- Classify releases into channels, e.g. stable & nightly (Done in Exchange server but the provider panel is WIP.)
- Extract related providers from `README.md` (WIP)
- Publish share to exchange server (WIP)
- Filter by provider information (WIP)
- Provider panel to show detail (WIP)

## Supported Providers

- Git
- GitHub
- NPM
- Docker
- YouTube

## Screenshot

![screenshot1](img/screenshot1.png)

## Future Plan

- Support more providers, e.g. Helm, PyPI, Facebook, Twitter, etc.
- Browser extensions
- Mobile apps
- Data import and export
- Tag with value
- Custom tag color
- Tag based sorting
- Webhook, client-side only
- Custom styling of share
- Preset : to show bookmarks with predefined conditions

## Getting started

#### Docker

```shell
docker run -d \
  -p 3399:3399 \
  --name pinmonl \
  pinmonl/pinmonl
```

If you would like to self host the Exchange server.

```shell
docker run -d \
  -p 8080:8080 \
  -e PINMONL_DB_DRIVER=sqlite3 \
  -e PINMONL_DB_DSN=exchange.db \
  --name pinmonl-exchange \
  pinmonl/exchange

docker run -d \
  -p 3399:3399 \
  -e PINMONL_EXCHANGE_ADDRESS=http://localhost:8080 \
  --name pinmonl \
  pinmonl/pinmonl
```

## Notes

1. By default, the bookmark listing page is showing only non-tagged item.
2. Scraper of DockerHub is not working.

## Key bindings

| Key   | Action          |
| ----- | --------------- |
| `/`   | Search          |
| `Esc` | Blur search box |

## License

MIT License.
