FROM golang:1.22.3-bullseye AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

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

# Install bash and required libraries
RUN apk add --no-cache bash libc6-compat

# Copy the app binary from the build stage
COPY --from=build /app/main /main

# Copy the .env file
COPY --from=build /app/pkg/envs/.env ./pkg/envs/.env

# Copy migrate binary
COPY --from=build /go/bin/migrate /usr/local/bin/migrate

# Copy migration files
COPY --from=build /app/internal/module/mail/migrations /migrations

# Copy template files
COPY --from=build /app/pkg/templates /pkg/templates

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]