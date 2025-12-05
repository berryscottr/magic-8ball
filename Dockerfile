# --- Go Build Stage ---
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o magic-8ball ./main.go

# --- Python Runtime Stage ---
FROM python:3.13-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    git wget curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN pip install --no-cache-dir playwright pandas \
    && playwright install --with-deps

COPY --from=builder /app/magic-8ball /app/magic-8ball
COPY --from=builder /app/data /app/data
COPY --from=builder /app/scripts /app/scripts

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./magic-8ball"]
