FROM golang:latest AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN cat go.mod && cat go.sum

RUN go mod download

COPY src/ ./src/

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./src/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/app .

EXPOSE 8080

CMD ["./app"]
