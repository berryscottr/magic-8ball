# ===== Builder stage =====
FROM golang:1.25

# Install minimal dependencies for Python 3.14 + Playwright + Go build
RUN apt-get update && apt-get install -y \
    python3.14 python3.14-venv python3.14-dev python3-pip \
    git wget curl build-essential ca-certificates \
    libnss3 libx11-xcb1 libxcomposite1 libxdamage1 libxrandr2 libasound2t64 \
    --no-install-recommends && rm -rf /var/lib/apt/lists/*

# Make python3 point to python3.14
RUN update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.14 1

WORKDIR /app

# Install Playwright system-wide
RUN python3 -m pip install --break-system-packages --upgrade pip
RUN python3 -m pip install --break-system-packages playwright
RUN playwright install

# Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy app source and build Go binary
COPY . ./
RUN go build -o magic-8ball ./main.go

# ===== Runtime stage =====
FROM python:3.14-slim

# Install only minimal runtime dependencies for Playwright
RUN apt-get update && apt-get install -y \
    libnss3 libx11-xcb1 libxcomposite1 libxdamage1 libxrandr2 libasound2t64 \
    wget ca-certificates \
    --no-install-recommends && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy Go binary + app files
COPY --from=0 /app/magic-8ball /app/magic-8ball
COPY --from=0 /app/data /app/data
COPY --from=0 /app/scripts /app/scripts

# Copy Python site-packages, Playwright CLI, and browser cache
COPY --from=0 /usr/local/lib/python3.14/dist-packages /usr/local/lib/python3.14/dist-packages
COPY --from=0 /usr/local/bin/playwright /usr/local/bin/playwright
COPY --from=0 /root/.cache/ms-playwright /root/.cache/ms-playwright

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./magic-8ball"]
