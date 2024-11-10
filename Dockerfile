# Dockerfile

# Stage 1: Build the Go application
FROM golang:1.23.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o web-crawler main.go

# Stage 2: Create a smaller final image
FROM alpine:latest
RUN apk add --no-cache libc6-compat

# Copy the built application and entrypoint script
COPY --from=builder /app/web-crawler /web-crawler
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Create the necessary directories
RUN mkdir -p /scraped-data /output

# Set the entrypoint to the script
ENTRYPOINT ["/entrypoint.sh"]
