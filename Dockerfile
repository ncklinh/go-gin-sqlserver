# Giai đoạn 1: Build app với Golang
FROM golang:1.21 AS builder

WORKDIR /app

# Copy và cài module
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ mã nguồn
COPY . .

# Build ra binary
RUN go build -o server .

# Giai đoạn 2: Chạy với distroless (tối giản, không có shell)
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary từ stage build
COPY --from=builder /app/server /server

# Mặc định port bạn expose (phải trùng với PORT trong env nếu có)
EXPOSE 8080

# Lệnh chạy app (không dùng shell)
ENTRYPOINT ["/server"]
