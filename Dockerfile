FROM golang:1.12.6-alpine3.9 AS builder

WORKDIR /go/src/github.com/kezmorris/gomessage
COPY main.go .
RUN go build .

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /go/src/github.com/kezmorris/gomessage/gomessage .

CMD ["./gomessage"]

