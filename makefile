export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt build

build: proxy_server proxy_client proxy_client2

fmt:
	go fmt ./...

proxy_server:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/proxy_server ./server

proxy_client:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/proxy_client ./client

proxy_client2:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/proxy_client2 ./client2

clean:
	rm -f ./bin/proxy_server
	rm -f ./bin/proxy_client