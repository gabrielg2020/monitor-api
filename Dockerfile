FROM golang:1.25-alpine AS builder
LABEL authors="gabriel"

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -ldflags="-w -s" -o app main.go

# Final stage
FROM alpine:latest

# Install SQLite3 runtime
RUN apk --no-cache add sqlite ca-certificates

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /build/app .

# Expose the port the app runs on
EXPOSE 8191

ENTRYPOINT ["./app"]