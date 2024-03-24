# Use the official Golang image as a base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /go/src/rapid

# Copy the go.mod and go.sum files to the container's workspace
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the source code to the container's workspace
COPY src .

# Enter shell for debugging
CMD ["sh"]
