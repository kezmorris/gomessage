FROM golang:1.12.9-alpine3.9 AS builder

RUN apk add git

ENV GOPATH=/go/
ENV GO111MODULE=on

COPY go.mod /go/src/github.com/kezmorris/gomessage/go.mod
COPY go.sum /go/src/github.com/kezmorris/gomessage/go.sum
WORKDIR /go/src/github.com/kezmorris/gomessage/

RUN go get

COPY operator/* /go/src/github.com/kezmorris/gomessage/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o /go/bin/gomessage

#TODO: upx up in here?
FROM alpine:3.9

COPY --from=builder /go/bin/gomessage .

CMD ["./gomessage"]
