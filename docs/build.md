# Build & Development

The project is built using [Go](https://golang.org/) 1.13+.

[![CircleCI](https://circleci.com/gh/cchantep/nuggan.svg?style=svg)](https://circleci.com/gh/cchantep/nuggan)

## Prerequisites

Install libvips:

- **macOS**: `port install vips` (using [MacPorts](https://www.macports.org/))
- **Linux**: Follow [libvips installation guide](https://libvips.github.io/libvips/install.html)

## Build

Build the `nuggan` executable from source:

```sh
go build
```

## Tests

Run the full test suite:

```sh
go test -v nuggan
```

## SSE Control (Apple Silicon & Non-SIMD Platforms)

The build system uses a compiler wrapper script (`scripts/cc-sse.sh`) to automatically handle SSE (SIMD) support based on your platform.

### Enable the SSE Wrapper

Set the `CC` environment variable to the absolute path of the wrapper script before building:

```sh
CC="$PWD/scripts/cc-sse.sh" go build
CC="$PWD/scripts/cc-sse.sh" go test -v nuggan
```

### How It Works

The wrapper automatically detects your platform:

- **ARM64 (Apple Silicon, etc.)** — SSE is automatically disabled since ARM64 doesn't support x86 SIMD instructions.
- **x86_64 (Intel/AMD, Linux CI)** — SSE is automatically enabled for optimal performance.

### Manual Override

Explicitly control SSE support using the `NUGGAN_USE_SSE` environment variable:

- `NUGGAN_USE_SSE=0` (or `false`, `no`, `off`) — Explicitly disable SSE
- `NUGGAN_USE_SSE=1` (or `true`, `yes`, `on`) — Explicitly enable SSE (may fail on ARM64)

### Build Examples

**ARM64 (SSE auto-disabled):**

```sh
CC="$PWD/scripts/cc-sse.sh" go build
```

**Explicitly disable SSE:**

```sh
NUGGAN_USE_SSE=0 CC="$PWD/scripts/cc-sse.sh" go build
```

**Run tests with SSE wrapper:**

```sh
CC="$PWD/scripts/cc-sse.sh" go test -v nuggan
```

### Note on Image Output Differences

Image output differs slightly when SSE is disabled due to different quantization algorithms. This is expected and normal. Output remains deterministic and valid across builds.

## Testing on Apple Silicon (Docker)

If you're on an Apple Silicon host or your local libvips setup differs from CI requirements, run tests in Docker using the CI environment:

```sh
docker run --rm --platform linux/amd64 \
  -v "$PWD":/go/src/github.com/cchantep/nuggan \
  -w /go/src/github.com/cchantep/nuggan \
  cchantep/golang:1.13-vips \
  /bin/sh -lc 'export PATH=/usr/local/go/bin:$PATH; go test -v nuggan'
```

This runs the exact test environment used in CI with x86_64 architecture and all dependencies properly configured.
