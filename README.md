# Pinmonl

## What is Pinmonl
A bookmark manager for developers:
- Shows package status and statistics, e.g. GitHub star, list of versions/tags.
- Share a list of bookmarks, which may be an **awesome** list replacement.
- Bookmark tagging
- Tracking free and Ad. free

Pinmonl is still in the development stage.

## Future Plan
Functional features must be part of the OSS and always be tracking & ad. free.

##### To-do
- Test cases
- CI, planned to use Drone
- Error handling
- Dark mode

##### Stage 1. Integrate with the API of major Git providers
- GitLab
- GitTea
- Gogs
- BitBucket

##### Stage 2. Integrate other service providers
- Helm
- Docker
- NPM
- Archlinux package
- and more...

##### Stage 3. Feature updates
- Single user mode
- Tag styling
- Shortcut key
- Webhook, and not going to use email.
- Share statistical data from self-hosted instance to Pinmonl server (There would be an option to opt-in/out).
- Pull request to the sharing list.
- Project derived bookmarking
- Mark on the release version (which may be useful to remember your favourite **STABLE** version)

##### Clients (to be developed)
- Mobile apps
- Browser extension
- CLI

## Get started

##### 1. Configuration
by `config.yaml`
```yaml
db:
  # driver: sqlite3
  # dsn: file:pinmonl.db?cache=shared
  # - OR -
  # driver: postgres
  # dsn: postgres://pinmonl:pinmonl@pg:5432/pinmonl?sslmode=disable
  # - OR -
  # driver: mysql
  # dsn: pinmonl:pinmonl@tcp(mysql:3306)/pinmonl
github:
  token: YourGitHubToken
cookie:
  hashkey: hash_with_at_least_32_length
  blockkey: hash_with_at_least_32_length
```

or use environment
```ini
DB_DRIVER=sqlite3
DB_DSN=file:pinmonl.db?cache=shared
GITHUB_TOKEN=YourGitHubToken
COOKIE_HASHKEY=hash_with_at_least_32_length
COOKIE_BLOCKKEY=hash_with_at_least_32_length
```

##### 2. Setup migration table
`pinmonl migration install`

##### 3. Run migrations
`pinmonl migration up`

##### 4. Start server
`pinmonl server`

## Support the project
[Patreon](https://patreon.com/pinmonl)

Pull requests are welcome.

## License
MIT License.
