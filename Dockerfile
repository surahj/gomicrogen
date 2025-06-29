# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gomicrogen .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S gomicrogen && \
    adduser -u 1001 -S gomicrogen -G gomicrogen

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/gomicrogen .

# Copy templates directory
COPY --from=builder /app/templates ./templates

# Change ownership to non-root user
RUN chown -R gomicrogen:gomicrogen /app

# Switch to non-root user
USER gomicrogen

# Expose port (if needed for any web interface)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./gomicrogen"] 