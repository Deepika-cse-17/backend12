# -------- BUILD STAGE --------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git (required for fetching modules)
RUN apk add --no-cache git

# Copy go files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -o server .

# -------- RUN STAGE --------
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 8081

# Run app
CMD ["./server"]