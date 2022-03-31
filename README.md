# Pad

[![build](https://github.com/dotzero/pad/actions/workflows/ci-build.yml/badge.svg)](https://github.com/dotzero/pad/actions/workflows/ci-build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotzero/pad)](https://goreportcard.com/report/github.com/dotzero/pad)
[![Docker Automated build](https://img.shields.io/docker/automated/jrottenberg/ffmpeg.svg)](https://hub.docker.com/r/dotzero/pad/)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dotzero/pad/blob/master/LICENSE)

Pad is a standalone version of Cloud Notepad. Allows you to share any text data via unique links.

Data stored in [boltdb](https://github.com/etcd-io/bbolt) (embedded key/value database) files under `PAD_DB_PATH`.

![](https://raw.githubusercontent.com/dotzero/pad/master/web/images/screenshot.png)

## Running container in Docker

```bash
docker run -d --rm --name pad \
    -p "8080:8080" \
    -e PAD_SECRET=random_salt_here \
    -e PAD_PORT=8080 \
    dotzero/pad
```

### Running container with Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: "3"
services:
  pad:
    build: .
    container_name: pad
    restart: always
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

```
git clone https://github.com/dotzero/pad
cd pad
go run .
```

### Command line options

```bash
Usage:
  pad [OPTIONS]

Application Options:
      --host=    listening address (default: 0.0.0.0) [$PAD_HOST]
      --port=    listening port (default: 8080) [$PAD_PORT]
      --db=      path to database files (default: db) [$PAD_DB_PATH]
      --path=    path to web assets (default: web) [$PAD_ASSETS_PATH]
      --secret=  the shared secret key used to generate ids [$PAD_SECRET]
      --verbose  verbose logging
  -v, --version  show the version number

Help Options:
  -h, --help     Show this help message
```

### Environment variables

* `PAD_HOST` (*default:* `0.0.0.0`) - listening address
* `PAD_PORT` (*default:* `8080`) - listening port
* `PAD_DB_PATH` (*default:* `$PWD/db`) - path to BoltDB database. It represents a consistent snapshot of your data
* `PAD_ASSETS_PATH` (*default:* `$PWD/web`) - path to web assets, templates and static files.
* `PAD_SECRET` (*default:* `empty`) - salt that using to generate hashids. **Strongly** recommend to replace with your own value

## License

http://www.opensource.org/licenses/mit-license.php
