# -------- BUILD STAGE --------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git (needed for go modules)
RUN apk add --no-cache git

# Copy only dependency files first (better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Now copy the rest of the source code
COPY . .

# Build static binary for Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# -------- RUNTIME STAGE --------
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 8081

# Run app
CMD ["./server"]