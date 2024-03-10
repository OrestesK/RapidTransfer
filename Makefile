BINARY_NAME=RapidTransfer
FLAGS=-ldflags "-s -w"

build:
	@go build ${FLAGS} -o ${BINARY_NAME} src/*.go

main:
	@go run src/main.go src/network.go src/parser.go $(ARG)

clean:
	@go clean
	@rm ${BINARY_NAME}

deps:
	@go get github.com/libp2p/go-libp2p
	@go get github.com/multiformats/go-multiaddr
	@go mod tidy

