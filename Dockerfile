# ===== Builder stage =====
FROM golang:1.25

# Install Python + pip + Playwright dependencies
RUN apt-get update && apt-get install -y \
    python3 python3-venv python3-pip \
    git wget curl build-essential ca-certificates \
    --no-install-recommends && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Install Playwright system-wide
RUN python3 -m pip install --break-system-packages playwright
RUN playwright install --with-deps

# Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy app and build Go binary
COPY . ./
RUN go build -o magic-8ball ./main.go

# ===== Runtime stage =====
FROM python:3.13-slim

WORKDIR /app

# Install Playwright runtime dependencies
RUN apt-get update && apt-get install -y \
    wget ca-certificates \
    --no-install-recommends && rm -rf /var/lib/apt/lists/*

# Copy Go binary + app files
COPY --from=0 /app/magic-8ball /app/magic-8ball
COPY --from=0 /app/data /app/data
COPY --from=0 /app/scripts /app/scripts

# Copy Python site-packages and Playwright CLI + browser cache
COPY --from=0 /usr/local/lib/python3.*/dist-packages /usr/local/lib/python3.*/dist-packages
COPY --from=0 /usr/local/bin/playwright /usr/local/bin/playwright
COPY --from=0 /root/.cache/ms-playwright /root/.cache/ms-playwright

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./magic-8ball"]
