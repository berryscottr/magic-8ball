# ===== Builder stage =====
FROM golang:1.25

# Install Go build deps + Python + Playwright dependencies
RUN apt-get update && apt-get install -y \
    git wget curl python3 python3-venv python3-pip \
    libnss3 libx11-xcb1 libxcomposite1 libxdamage1 libxrandr2 libasound2t64 \
    --no-install-recommends && rm -rf /var/lib/apt/lists/*

# Set up Python virtual environment
RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Install Playwright
RUN pip install --upgrade pip && pip install playwright && playwright install

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o magic-8ball ./main.go

# ===== Runtime stage =====
FROM ubuntu:24.04

# Install minimal Python + dependencies for Playwright
RUN apt-get update && apt-get install -y \
    python3 python3-venv python3-pip \
    libnss3 libx11-xcb1 libxcomposite1 libxdamage1 libxrandr2 libasound2t64 \
    --no-install-recommends && rm -rf /var/lib/apt/lists/*

# Copy Python venv + Go binary from builder
COPY --from=0 /opt/venv /opt/venv
COPY --from=0 /app/magic-8ball /app/magic-8ball

ENV PATH="/opt/venv/bin:$PATH"
WORKDIR /app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./magic-8ball"]
