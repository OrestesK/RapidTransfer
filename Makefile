EXECUTABLE=Rapid

build:
	go mod download

windows:
	env GOOS=windows GOARCH=amd64 go build -o $$GOPATH/bin/$(EXECUTABLE).exe src/*.go

linux:
	env GOOS=linux GOARCH=amd64 go build -v -o $$GOPATH/bin/$(EXECUTABLE) src/*.go

darwin:
	env GOOS=darwin GOARCH=amd64 go build -v -o $$GOPATH/bin/$(EXECUTABLE) src/*.go