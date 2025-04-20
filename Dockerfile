# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ride-matcher ./cmd/ride-matcher

# Final stage
FROM alpine:latest

# Add timezone support
RUN apk --no-cache add tzdata

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/ride-matcher .

# Set ownership to the non-root user
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD ps aux | grep ride-matcher || exit 1

# Set the entrypoint
ENTRYPOINT ["./ride-matcher"]
