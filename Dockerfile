FROM golang:1.22-alpine

WORKDIR /app

# Copy everything
COPY . .

# Fix dependencies inside container
RUN go mod tidy

# Build binary
RUN go build -o server .

EXPOSE 8081

CMD ["./server"]