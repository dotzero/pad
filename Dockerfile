FROM golang:1.17-alpine AS build-env

# package dependencies
RUN apk add --update --no-cache build-base make git libc-dev

WORKDIR /src

COPY . .
RUN make build \
    && /go/bin/pad --version

FROM alpine:3.15

WORKDIR /app

# copy artefacts
COPY --from=build-env /go/bin/pad /app
COPY --from=build-env /src/web/ /app/web/

ENV PAD_DB_PATH "./db"
ENV PAD_SECRET "true_random_salt"
ENV PAD_HOST "0.0.0.0"
ENV PAD_PORT "8080"

EXPOSE ${PAD_PORT}

ENTRYPOINT ["/app/pad"]
