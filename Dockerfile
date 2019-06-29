FROM golang:1.12.6 AS builder

WORKDIR /go/src/github.com/kezmorris/gomessage
COPY main.go .
RUN go build .

FROM alpine:latest

COPY --from=builder /go/src/github.com/kezmorris/gomessage/gomessage /usr/bin/gomessage

CMD ["gomessage"]

