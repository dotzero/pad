# Pad

[![Go Report Card](https://goreportcard.com/badge/github.com/dotzero/pad)](https://goreportcard.com/report/github.com/dotzero/pad)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dotzero/pad/blob/master/LICENSE)

Pad is a standalone version of cloud notepad. Allows to share any text data by unique links.
Written in Go and Bolt as embedded key/value database.

![](https://raw.githubusercontent.com/dotzero/pad/master/static/images/screenshot.png)

## Installation

```bash
git clone https://github.com/dotzero/pad.git
cd pad
```

### Running container

```bash
> docker build -t dotzero/pad .

> docker run -d --name pad_app \
	-p  "8080:8080" \
	-e PAD_SALT=random_salt_here \
	-e PAD_PORT=8080 \
	dotzero/pad
```

### Running container with `docker-compose`

Create a `docker-compose.yml` file:

```
version: "3"
services:
  pad:
    build: .
    container_name: pad-app
    restart: always
    ports:
      - "8080:8080"
    environment:
        PAD_SALT: random_salt_here
        PAD_PORT: 8080
```

Run `docker-compose up -d`.

## Environment variables

### PAD_DB

* *default:* `pad.db`

Name of BoltDB database filename. It represents a consistent snapshot of your data.

### PAD_SALT

* *default:* `true_random_salt`

Salt that using to generate hashids. Strongly recommend to replace with your own value.

### PAD_PORT

* *default:* `8080`

This port **must** match the port that is exposed via Docker.

## License

http://www.opensource.org/licenses/mit-license.php
