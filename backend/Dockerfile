FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build both binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate/

FROM alpine:latest

# Install postgresql-client for health checks
RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /app

# Copy both binaries
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

# Create migrations directory and copy migrations
RUN mkdir -p /app/migrations
COPY --from=builder /app/migrations/*.sql /app/migrations/

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

# Command to run the executable
CMD ["./main"]