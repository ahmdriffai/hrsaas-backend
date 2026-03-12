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

WORKDIR /app

# Copy executable dari stage build
COPY --from=builder /app/build/main .

# Expose port yang digunakan oleh aplikasi
EXPOSE 9000

# Jalankan aplikasi
CMD ["./main"]