FROM golang:1.22.3-bullseye AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

COPY .. .

# Compile application during build rather than at runtime
# Add flags to statically link binary
RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o main ./cmd/main.go

####
FROM alpine:latest AS app

# Set gin mode
ENV GIN_MODE=release

WORKDIR /

# Copy the app binary from the build stage
COPY --from=build /app/main /main

# Copy the .env file
COPY --from=build /app/pkg/envs/.env ./pkg/envs/.env

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]