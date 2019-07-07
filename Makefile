BINARY=gomessage
IMAGE_NAME=morriski/gomessage
TAG_NAME=latest-dev

.DEFAULT_GOAL: $(BINARY)

#There is no windows
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
		os := darwin
else
		os := linux
endif

$(BINARY): clean
	mkdir -p bin
	env CGO_ENABLED=0 GOOS=${os} GOARCH=amd64 go build -ldflags="-s -w" -a -o bin/${BINARY}-server server/main.go
	env CGO_ENABLED=0 GOOS=${os} GOARCH=amd64 go build -ldflags="-s -w" -a -o bin/${BINARY}-client client/main.go

skinny: clean
	env CGO_ENABLED=0 GOOS=${os} GOARCH=amd64 go build -ldflags="-s -w" -a -o bin/${BINARY}-server server/main.go
	env CGO_ENABLED=0 GOOS=${os} GOARCH=amd64 go build -ldflags="-s -w" -a -o bin/${BINARY}-client client/main.go
	@upx --brute bin/${BINARY}-server
	@upx --brute bin/${BINARY}-client

clean: 
	go clean server/main.go
	rm -rf ${BINARY}-server
	go clean client/main.go
	rm -rf ${BINARY}-client

docker: 
	docker build . -f Dockerfile -t $(IMAGE_NAME):$(TAG_NAME)
	docker push $(IMAGE_NAME):$(TAG_NAME)
