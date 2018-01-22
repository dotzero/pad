FROM golang:1.9-alpine AS build-env

# Package dependencies
RUN apk add --update --no-cache make git libc-dev

# install dependency tool
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/dotzero/pad
COPY . .

RUN dep ensure -v \
    && make build \
    && pad -version \
    && which pad

FROM alpine:3.7

# Copy html to nginx
COPY --from=build-env /go/bin/pad /usr/bin/caddy

EXPOSE 8080

ENTRYPOINT ["/usr/bin/caddy"]
