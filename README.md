# Nuggan

[![CircleCI](https://circleci.com/gh/cchantep/nuggan.svg?style=svg)](https://circleci.com/gh/cchantep/nuggan)

Image optimizer HTTP micro-service & utilities based on [libvips](https://libvips.github.io/libvips/).

## Motivation

Nuggan provides a lightweight image proxy for on-the-fly cropping and resizing. Unlike many alternatives, it works with any existing web-hosted images without requiring additional infrastructure or DNS control.

## Quick Start

Run as a standalone HTTP service:

```sh
./nuggan -server ':8080' -server-config server.conf
```

## Documentation

- **[Usage Guide](./docs/usage.md)** — How to run (standalone, Docker, configuration)
- **[API Reference](./docs/api.md)** — Request format, parameters, and examples
- **[Build & Development](./docs/build.md)** — Building from source, testing, SSE control
- **[Deployment](./docs/deployment.md)** — Production deployment options (standalone, serverless, AWS Lambda)
- **[Dependencies](./docs/dependencies.md)** — Dependency management with Dependabot

For a complete overview, see [docs/index.md](./docs/index.md).

## Build Locally

Prerequisites: [libvips](https://libvips.github.io/libvips/install.html) (macOS: `port install vips`)

```sh
go build
go test -v nuggan
```

See [Build & Development](./docs/build.md) for full build instructions including SSE control for Apple Silicon.

## License

See [LICENSE](LICENSE) file.