# Pad

[![build](https://github.com/dotzero/pad/actions/workflows/ci-build.yml/badge.svg)](https://github.com/dotzero/pad/actions/workflows/ci-build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotzero/pad)](https://goreportcard.com/report/github.com/dotzero/pad)
[![Docker Automated build](https://img.shields.io/docker/automated/jrottenberg/ffmpeg.svg)](https://hub.docker.com/r/dotzero/pad/)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dotzero/pad/blob/master/LICENSE)

Pad is a standalone version of cloud notepad. Allows to share any text data by unique links.
Written in Go and Bolt as embedded key/value database.

![](https://raw.githubusercontent.com/dotzero/pad/master/web/images/screenshot.png)

## Running container

```bash
> docker run -d --rm --name pad_app \
    -p "8080:8080" \
    -e PAD_SECRET=random_salt_here \
    -e PAD_PORT=8080 \
    dotzero/pad
```

### Running container with `docker-compose`

Create a `docker-compose.yml` file:

```
version: "3"
services:
  pad:
    image: dotzero/pad
    container_name: pad_app
    restart: always
    ports:
      - "8080:8080"
    environment:
      PAD_DB_PATH: /app/db
      PAD_SECRET: random_salt_here
      PAD_PORT: 8080
    volumes:
      - ./db:/app/db
```

Run `docker-compose up -d`, wait for it to initialize completely, and visit `http://localhost:8080`

### Build container

```bash
> docker build -t dotzero/pad .
```

## Usage

```
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

## Environment variables

### PAD_HOST

* *default:* `0.0.0.0`

Listening address

### PAD_PORT

* *default:* `8080`

Listening port.

### PAD_DB_PATH

* *default:* `$PWD/db`

Path to BoltDB database. It represents a consistent snapshot of your data.

### PAD_ASSETS_PATH

* *default:* `$PWD/web`

Path to web assets, templates and static files.

### PAD_SECRET

* *default:* `empty`

Salt that using to generate hashids. Strongly recommend to replace with your own value.

## License

http://www.opensource.org/licenses/mit-license.php
