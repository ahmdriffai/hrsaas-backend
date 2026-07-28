# Stage 1: Build
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum (jika ada)
COPY go.mod ./

# Download dependencies (cache layer ini untuk efisiensi build)
RUN go mod download

COPY . .

# Build aplikasi
RUN go build -o build/main cmd/web/main.go 

# Stage 2: Minimal runtime
FROM alpine:latest

# tzdata is required for Go's time.Local to resolve named zones (e.g. TZ=Asia/Jakarta);
# alpine ships with no timezone database by default and Go silently falls back to UTC without it.
RUN apk add --no-cache tzdata

ENV TZ=Asia/Jakarta

WORKDIR /app

# Copy executable dari stage build
COPY --from=builder /app/build/main .

# Expose port yang digunakan oleh aplikasi
EXPOSE 9000

# Jalankan aplikasi
CMD ["./main"]