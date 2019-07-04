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
	env CGO_ENABLED=0 GOOS=${os} GOARCH=amd64 go build -ldflags="-s -w" -a -o ${BINARY}

clean: 
	go clean
	rm -rf ${BINARY}

docker: 
	docker build . -f Dockerfile
	docker push $(IMAGE_NAME):$(TAG_NAME)
