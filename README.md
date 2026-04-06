# Pad

[![build](https://github.com/dotzero/pad/actions/workflows/ci.yml/badge.svg)](https://github.com/dotzero/pad/actions/workflows/ci.yml)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dotzero/pad/blob/master/LICENSE)

Pad is a standalone version of Cloud Notepad that lets you share text through unique links.

Data is stored in [BoltDB](https://github.com/etcd-io/bbolt) files under `BOLT_PATH`.

![](./images/screenshot.png)

## Run with Docker

```bash
docker run -d --rm --name pad -p "8080:8080" dotzero/pad
```

### Run with Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: "3"
services:
  pad:
    image: ghcr.io/dotzero/pad:latest
    container_name: pad
    restart: always
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "5"
    ports:
      - "8080:8080"
    environment:
      PAD_HOST: "0.0.0.0"
      PAD_PORT: "8080"
      PAD_SECRET: random_salt_here
    volumes:
      - pad_db:/app/db

volumes:
  pad_db:
```

Run `docker compose up -d`, wait for the container to start, then open `http://localhost:8080`.

### Build the container image

```bash
docker build -t dotzero/pad .
```

## Run locally

```bash
git clone https://github.com/dotzero/pad
cd pad
go run .
```

### Command line options

```bash
Usage:
  pad [OPTIONS]

Application Options:
      --host=        listening address (default: 0.0.0.0) [$PAD_HOST]
      --port=        listening port (default: 8080) [$PAD_PORT]
      --bolt-path=   parent directory for the bolt files (default: ./var) [$BOLT_PATH]
      --secret=      the shared secret key used to generate ids [$PAD_SECRET]
      --static-path= path to website assets (default: ./static) [$STATIC_PATH]
      --tpl-path=    path to templates files (default: ./templates) [$TPL_PATH]
      --tpl-ext=     templates files extensions (default: .html) [$TPL_EXT]
      --verbose      verbose logging
  -v, --version      show the version number

Help Options:
  -h, --help         Show this help message
```

### Environment variables

- `PAD_HOST` (default: `0.0.0.0`) - Listening address.
- `PAD_PORT` (default: `8080`) - Listening port.
- `BOLT_PATH` (default: `./var`) - Path to the BoltDB data directory.
- `PAD_SECRET` (default: empty) - Salt used to generate hashids. Replace it with your own secret in production.
- `STATIC_PATH` (default: `./static`) - Path to static assets.
- `TPL_PATH` (default: `./templates`) - Path to template files.
- `TPL_EXT` (default: `.html`) - Template file extension.

## License

[MIT](http://www.opensource.org/licenses/mit-license.php)
