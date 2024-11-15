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

# Copy the built application
COPY --from=builder /app/web-crawler /web-crawler

# Set the command to run the application
CMD ["/web-crawler"]