# Stage 1: Builder
FROM golang:1.23-alpine as builder

# Set build arguments and environment variables
ARG BOT_TOKEN
ENV BOT_TOKEN=$BOT_TOKEN

# Install dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . ./

# Build the Go application
RUN go build -o magic-8ball ./main.go

# Stage 2: Runtime
FROM alpine:latest as runtime

# Install minimal dependencies for runtime
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/magic-8ball .

# Expose the application port
EXPOSE 8080

# Add a health check
# HEALTHCHECK --interval=30s --timeout=10s \
#   CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# Start the application
CMD ["./magic-8ball"]
