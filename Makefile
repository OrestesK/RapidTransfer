BINARY_NAME= bin/RapidTransfer
FLAGS=-ldflags "-s -w"

build:
	@go build ${FLAGS} -o ${BINARY_NAME} src/main.go

server: build
	@./${BINARY_NAME}

client: build
	@./${BINARY_NAME} -p $(key)

clean:
	@go clean
	@rm ${BINARY_NAME}

deps:
	@go get github.com/libp2p/go-libp2p
	@go get github.com/multiformats/go-multiaddr
	@go mod tidy

