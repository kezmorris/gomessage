FROM golang:1.12.6-alpine3.9 AS builder

ENV GOPATH=/go/
WORKDIR /go/src/github.com/kezmorris/gomessage/
COPY operator/* ./
RUN go build .
#TODO: upx up in here?
FROM alpine:latest

COPY --from=builder /go/src/github.com/kezmorris/gomessage/gomessage .

CMD ["./gomessage"]
