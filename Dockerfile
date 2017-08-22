FROM golang:1.8-alpine

WORKDIR /go/src/github.com/dotzero/pad
COPY . .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]
