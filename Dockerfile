## Build stage
FROM golang:alpine AS build

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
COPY ./vendor ./vendor
RUN go mod download

## Deployment stage
FROM golang:alpine

WORKDIR /app

COPY --from=build /app/go.mod .
COPY --from=build /app/go.sum .
COPY --from=build /app/vendor .

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go build -o ./bin/sqlite ./cmd/sqlite/sqlite_server.go

EXPOSE 8080/tcp

ENTRYPOINT [ "/app/bin/sqlite"]
