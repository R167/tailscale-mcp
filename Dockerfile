FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -o tailscale-mcp .

# Final stage
FROM scratch

LABEL org.opencontainers.image.source="https://github.com/R167/tailscale-mcp"
LABEL org.opencontainers.image.description="Model Context Protocol server for Tailscale"
LABEL org.opencontainers.image.licenses="MIT"

# Copy ca-certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/tailscale-mcp /tailscale-mcp

# Expose the default port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/tailscale-mcp"]