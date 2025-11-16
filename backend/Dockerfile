FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Install dependencies including postgresql-client for pg_isready
RUN apk add --no-cache postgresql-client

# Install Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and migrations
COPY . .
COPY migrations/ ./migrations/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate/

FROM alpine:latest

# Install postgresql-client for health checks
RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /root/

# Copy the pre-built binary files AND migrations
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations/

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]