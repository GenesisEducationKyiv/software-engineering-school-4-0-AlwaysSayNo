# syntax=docker/dockerfile:1
FROM golang:1.22.3-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Install dependencies
RUN apk update && apk add --no-cache git build-base

# Copy go mod and sum files
COPY go.mod go.sum ./

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN --mount=type=cache,target=/go/pkg/mod \
             --mount=type=cache,target=/root/.cache/go-build \
             go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build \
  -ldflags="-linkmode external -extldflags -static" \

  -tags netgo \
  -o main ./cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]