FROM golang:1.8-alpine

# Package dependencies
RUN apk add --update --no-cache make git libc-dev

# install dependency tool
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/dotzero/pad
COPY . .

EXPOSE 8080

RUN dep ensure -v \
    && make build \
    && pad -version

CMD ["go-wrapper", "run"]
