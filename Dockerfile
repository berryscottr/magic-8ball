# ===== Builder stage =====
FROM golang:1.25 as builder

# Install system dependencies for Python + Playwright + Chromium
RUN apt-get update && apt-get install -y \
    git wget curl python3 python3-venv python3-pip \
    chromium chromium-driver \
    libnss3 libx11-xcb1 libxcomposite1 libxdamage1 \
    libxrandr2 libasound2 --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Set up Python virtual environment for Playwright
RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

RUN pip install --upgrade pip
RUN pip install playwright
RUN playwright install --with-deps chromium

# Build Go app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o magic-8ball ./main.go

# ===== Runtime stage =====
FROM debian:slim

# Install runtime dependencies for Python + Chromium
RUN apt-get update && apt-get install -y \
    python3 python3-venv python3-pip chromium chromium-driver \
    libnss3 libx11-xcb1 libxcomposite1 libxdamage1 \
    libxrandr2 libasound2 --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Copy Go binary and Python virtual environment from builder
COPY --from=builder /app/magic-8ball /app/magic-8ball
COPY --from=builder /opt/venv /opt/venv

ENV PATH="/opt/venv/bin:$PATH"
WORKDIR /app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./magic-8ball"]
