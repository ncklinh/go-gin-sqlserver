# Stage 1: Build app with Golang
FROM golang:1.21 AS builder

WORKDIR /app

# Copy and install modules
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build binary
RUN go build -o server .

# Stage 2: Run with distroless (minimal, no shell)
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/server /server

# Default port to expose (must match PORT in env if any)
EXPOSE 8080

# Command to run app (no shell)
ENTRYPOINT ["/server"]
