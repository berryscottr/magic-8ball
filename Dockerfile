# Use Debian-slim base for better Python compatibility
FROM golang:1.25-bullseye as builder

# Install dependencies for Go build and Python/Playwright
RUN apt-get update && apt-get install -y \
    git wget curl python3 python3-venv python3-pip \
    chromium chromium-driver fonts-liberation \
    libnss3 libx11-xcb1 libxcomposite1 libxcursor1 \
    libxdamage1 libxi6 libxtst6 libxrandr2 libasound2 \
    libatk1.0-0 libcups2 libxss1 libglib2.0-0 \
    libpangocairo-1.0-0 libpango-1.0-0 libfontconfig1 \
    libfreetype6 libharfbuzz0b --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Set up Python venv and install Playwright
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

# Runtime image
FROM debian:bullseye-slim

# Install runtime dependencies for Python and Chromium
RUN apt-get update && apt-get install -y \
    python3 python3-venv python3-pip chromium chromium-driver \
    libnss3 libx11-xcb1 libxcomposite1 libxcursor1 \
    libxdamage1 libxi6 libxtst6 libxrandr2 libasound2 \
    libatk1.0-0 libcups2 libxss1 libglib2.0-0 \
    libpangocairo-1.0-0 libpango-1.0-0 libfontconfig1 \
    libfreetype6 libharfbuzz0b wget curl git --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Copy Go binary and Python venv from builder
COPY --from=builder /app/magic-8ball /app/magic-8ball
COPY --from=builder /opt/venv /opt/venv

ENV PATH="/opt/venv/bin:$PATH"
WORKDIR /app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./magic-8ball"]
