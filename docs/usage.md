# Usage Guide

## Standalone Service

Run Nuggan as an HTTP service on port `8080`:

```sh
./nuggan -server ':8080' -server-config server.conf
```

Then images available through configured base URLs are served at URLs like:

```
http://localhost:8080/optimg/0/0/-/-/-/-/-/_2_L3BvcHRvY2F0X3YyLnBuZw==/image.png
```

The image is first cropped (if crop parameters are specified), then resized (if resize parameters are specified). The original image is resolved using the base64-encoded `base64Ref` parameter.

For detailed URL format and parameter reference, see the [API Reference](./api.md).

## Docker

### Run from Docker Image

Nuggan can be easily started as a Docker container:

```sh
docker run --rm -P cchantep/nuggan:amazonlinux1
```

### Override Configuration

Override the default configuration with a volume:

```sh
docker run -v /tmp/custom.conf:/root/server.conf \
  --rm -P cchantep/nuggan:amazonlinux1
```

### Build Docker Image Locally

Build the Docker image locally:

```sh
./scripts/docker-build.sh
```

## Configuration

Create a server configuration file (e.g., `server.conf`):

```
groupedBaseUrls = [
  [
    "https://upload.wikimedia.org/wikipedia/commons"
  ],
  [
    "https://cdn0.iconfinder.com/data/icons",
    "https://cdn1.iconfinder.com/data/icons"
  ],
  [
    "https://octodex.github.com/images"
  ]
]
routePrefix = "optimg"
strict = true
cacheControl = "max-age=7200, s-maxage=21600"
```

### Configuration Fields

- **`groupedBaseUrls`**: A list of groups of URLs. Each group specifies base URLs corresponding to a same image source. Used in strict mode to validate image references.
- **`routePrefix`**: The prefix for the HTTP image API (default: `optimg`). This appears in all request URLs.
- **`strict`**: Strict mode (default: `false`). When enabled, only images from the configured sources in `groupedBaseUrls` can be requested. In strict mode, image references must follow the format `_{groupIndex}_{base64ImagePath}`.
- **`cacheControl`**: Optional [`Cache-Control`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control) response header. Example: `"max-age=7200, s-maxage=21600"`.

## Utilities

### Encode Image URLs

Encode an image URL to a `base64Ref` accepted by the service according to your server configuration:

```sh
./nuggan -server-config server.conf -encode-url https://octodex.github.com/images/poptocat_v2.png
```

Output:

```
2020/01/12 16:06:07 Encode 'https://octodex.github.com/images/poptocat_v2.png':

	_2_L3BvcHRvY2F0X3YyLnBuZw==
...
```

This encoded reference can then be used directly in image requests. The encoding depends on your `strict` mode setting and configured URL groups—see the [API Reference](./api.md) for details.
