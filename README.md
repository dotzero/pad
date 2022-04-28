# Pad

[![build](https://github.com/dotzero/pad/actions/workflows/ci-build.yml/badge.svg)](https://github.com/dotzero/pad/actions/workflows/ci-build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotzero/pad)](https://goreportcard.com/report/github.com/dotzero/pad)
[![Docker Automated build](https://img.shields.io/docker/automated/jrottenberg/ffmpeg.svg)](https://hub.docker.com/r/dotzero/pad/)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dotzero/pad/blob/master/LICENSE)

Pad is a standalone version of Cloud Notepad. Allows you to share any text data via unique links.

Data stored in [boltdb](https://github.com/etcd-io/bbolt) (embedded key/value database) files under `BOLT_PATH`.

![](https://raw.githubusercontent.com/dotzero/pad/master/static/images/screenshot.png)

## Running container in Docker

```bash
docker run -d --rm --name pad -p "8080:8080" dotzero/pad
```

### Running container with Docker Compose

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

Run `docker-compose up -d`, wait for it to initialize completely, and visit `http://localhost:8080`

### Build container

```bash
docker build -t dotzero/pad .
```

## How to run it locally

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

* `PAD_HOST` (*default:* `0.0.0.0`) - listening address
* `PAD_PORT` (*default:* `8080`) - listening port
* `BOLT_PATH` (*default:* `./var`) - path to BoltDB database (it represents a consistent snapshot of your data)
* `PAD_SECRET` (*default:* `empty`) - salt that using to generate hashids. **Strongly** recommend to replace with your own value
* `STATIC_PATH` (*default:* `./static`) - path to web assets
* `TPL_PATH` (*default:* `./templates`) - path to templates
* `TPL_EXT` (*default:* `.html`) - templates files extensions

## License

http://www.opensource.org/licenses/mit-license.php
