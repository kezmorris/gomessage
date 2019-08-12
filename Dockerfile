FROM golang:1.12.6-alpine3.9 AS builder

WORKDIR /go/src/github.com/kezmorris/gomessage
COPY server/main.go .
RUN go build .
#TODO: upx up in here?
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /go/src/github.com/kezmorris/gomessage/gomessage .

CMD ["./gomessage"]
EXPOSE 8001
