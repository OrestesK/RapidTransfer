EXECUTABLE=Rapid
SRC_DIR := Rapid
DEST_DIR := $(USERPROFILE)

# Moving the directory and account information into the users directory
prep:
	@echo "Moving directory $(SRC_DIR) to $(DEST_DIR)"
	mv $(SRC_DIR)\* $(DEST_DIR)
	@echo "Directory moved successfully."

# Downloads all dependencies
build:
	go mod download

# Creates binary for users on windows
windows:
	env GOOS=windows GOARCH=amd64 go build -o $$GOPATH/bin/$(EXECUTABLE).exe src/*.go
