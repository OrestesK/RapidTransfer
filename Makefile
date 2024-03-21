EXECUTABLE=Rapid
SRC_DIR := Rapid
DEST_DIR := $(USERPROFILE)
detected_OS := $(shell uname)

# Moving the directory and account information into the users directory
prep:
	@echo "Moving directory $(SRC_DIR) to $(DEST_DIR)"
	mv $(SRC_DIR)\* $(DEST_DIR)
	@echo "Directory moved successfully."

# Downloads all dependencies
deps:
	go mod download

build:
	@if [ "$(findstring MINGW64,$(detected_OS))" = "MINGW64" ]; then \
		env GOOS=windows GOARCH=amd64 go build -o $$GOPATH/bin/$(EXECUTABLE).exe src/*.go; \
	fi
	@if [ "$(detected_OS)" = "Darwin" ]; then \
		env GOOS=darwin GOARCH=amd64 go build -o $$(go env GOPATH)/bin/$(EXECUTABLE) src/*.go; \
	fi
	@if [ "$(detected_OS)" = "Linux" ]; then \
		env GOOS=linux GOARCH=amd64 go build -o $$(go env GOPATH)/bin/$(EXECUTABLE) src/*.go; \
	fi