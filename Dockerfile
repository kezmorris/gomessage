FROM golang:1.11.9-alpine3.9 AS builder

ENV GOPATH=/go/
COPY operator/* ./
WORKDIR /go/src/github.com/kezmorris/gomessage/




RUN go build .
#TODO: upx up in here?
FROM alpine3.9:latest

COPY --from=builder /go/src/github.com/kezmorris/gomessage/gomessage .

CMD ["./gomessage"]
