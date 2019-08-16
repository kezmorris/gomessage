FROM golang:1.12.9-alpine3.9 AS builder

RUN apk add git

ENV GOPATH=/go/
COPY operator/* /go/src/github.com/kezmorris/gomessage/
WORKDIR /go/src/github.com/kezmorris/gomessage/
RUN go get 
RUN go build .
#TODO: upx up in here?
FROM alpine:3.9

COPY --from=builder /go/src/github.com/kezmorris/gomessage/gomessage .

CMD ["./gomessage"]
