# Sentinel

- [Overview](#overview)
- [Architecture](#architecture)
- [Configuration](#configuration)
- [Main Components](#main-components)
  - [Main Application (`main.go`)](#main-application-maingo)
  - [Configuration Module (`pkg/config`)](#configuration-module-pkgconfig)
  - [Server Module (`pkg/server`)](#server-module-pkgserver)
- [API Usage](#api-usage)
- [Docker & Deployment](#docker--deployment)
- [Development & Dependencies](#development--dependencies)
- [CI/CD](#cicd)
- [Example Configuration](#example-configuration)

---

## Overview

**Sentinel** is a lightweight HTTP proxy for forwarding Sentry envelopes to upstream Sentry servers. It is designed to route requests from multiple projects, providing a secure and configurable way to forward error and event data to Sentry, with per-project authentication and routing.

## Architecture

Sentinel is composed of three main modules:
- **Main Application**: Entry point, configuration loading, and server startup.
- **Configuration Module**: Loads, validates, and provides access to runtime configuration.
- **Server Module**: HTTP server, request validation, and proxy logic.

## Configuration

Sentinel uses a JSON configuration file (default: `settings.json`). The configuration includes:
- List of projects and their Sentry credentials
- Server port and route
- Upstream Sentry URL
- Custom header for project selection
- Allowed origins for CORS (array of strings)
- `escapePlaceholder`: a placeholder value you can use in your Sentry configuration (e.g., DSN field). Sentinel will replace every occurrence of this value in incoming envelopes with the actual DSN for the target project before forwarding to Sentry. This allows you to avoid hardcoding DSNs in client configs and centralize routing logic.

See [Example Configuration](#example-configuration) for details.

## Main Components

### Main Application (`main.go`)
- Loads configuration using [Viper](https://github.com/spf13/viper) and the custom config loader.
- Sets default values and config search paths.
- Starts the HTTP server via `server.Run()`.
- Panics on configuration or server startup errors.

### Configuration Module (`pkg/config`)
- Defines the configuration schema (`Config`, `Server`, `Project`).
- Loads and parses the JSON config file.
- Validates:
  - At least one project is defined
  - Sentry URL is valid, not an IP, and uses HTTP/HTTPS
  - Server port is in range 1-65535
  - Route starts with `/` and does not end with `/`
  - `allowedOrigins` is a non-empty array
- Exposes the loaded config as a global `Instance`.

### Server Module (`pkg/server`)
- Uses [Gin](https://github.com/gin-gonic/gin) for HTTP routing and [gin-contrib/cors](https://github.com/gin-contrib/cors) for CORS.
- Registers a POST endpoint at the configured route.
- Validates the presence of a custom header (e.g., `Sec-Ner` or `x-Sentry-App-Selector`).
- Looks up the project by header value; rejects if not found.
- Reads the request body and forwards it to the upstream Sentry envelope endpoint, using the project's credentials.
- Replaces all occurrences of the configured `escapePlaceholder` value in the request body with the actual DSN for the target project, then forwards the modified envelope to the upstream Sentry endpoint using the project's credentials.
- Returns 200 on success, propagates Sentry errors otherwise.
- CORS is configured to allow only the origins specified in `allowedOrigins` from the configuration file.

## API Usage

- **POST** `{route}`
  - Header: `{header}` (e.g., `Sec-Ner`)
  - Body: Sentry envelope
  - Response: 200 OK on success, 400/500 on error

## Docker & Deployment

- Multi-stage Docker build (see `Dockerfile`):
  - Builds a static Go binary for `linux/amd64`
  - Final image is based on Alpine Linux
  - Entrypoint: `sentinel`
- Example usage:
  ```sh
  docker build -t sentinel .
  docker run -v $(pwd)/settings.json:/settings.json -p 3000:3000 sentinel
  ```

## Development & Dependencies

- Go 1.24+
- Key dependencies:
  - `github.com/gin-gonic/gin` (HTTP server)
  - `github.com/gin-contrib/cors` (CORS)
  - `github.com/spf13/viper` (config management)
- See `go.mod` for full dependency list.

## CI/CD

- GitHub Actions workflow: `.github/workflows/docker-image.yml`
  - Builds and pushes Docker images to GitHub Container Registry on push/PR to `master`.

## Example Configuration

### `settings.example.json`
```json
{
  "projects": {
    "appid1": {
      "sentryProjectId": "1",
      "sentryKey": "hex"
    },
    "appid2": {
      "sentryProjectId": "2",
      "sentryKey": "hex"
    }
  },
  "server": {
    "port": 3000,
    "route": "/tunnel"
  },
  "sentryUrl": "https://sentry.io",
  "escapePlaceholder": "https://0@sentry.io/0",
  "header": "x-Sentry-App-Selector",
  "allowedOrigins": [
    "https://example.com",
    "http://localhost:3000"
  ]
}
```

---

For more details, see the source files and comments.
