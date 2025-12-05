# ===== Builder stage =====
FROM golang:1.25

# Install Go build deps + Python + Playwright dependencies
RUN apt-get update && apt-get install -y \
    git wget curl python3 python3-pip python3-venv \
    libnss3 libx11-xcb1 libxcomposite1 libxdamage1 libxrandr2 libasound2t64 \
    --no-install-recommends && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Install Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy app source and build Go binary
COPY . ./
RUN go build -o magic-8ball ./main.go

# Install Python deps + Playwright in a relocatable folder
RUN mkdir /opt/pydeps
RUN pip install --upgrade pip
RUN pip install --target=/opt/pydeps playwright
RUN /opt/pydeps/bin/playwright install

# ===== Runtime stage =====
FROM ubuntu:24.04

# Install minimal Python + dependencies
RUN apt-get update && apt-get install -y \
    python3 python3-pip \
    libnss3 libx11-xcb1 libxcomposite1 libxdamage1 libxrandr2 libasound2t64 \
    wget \
    --no-install-recommends && rm -rf /var/lib/apt/lists/*

# Copy app files + Go binary + Python deps from builder
COPY --from=0 /app/data /app/data
COPY --from=0 /app/scripts /app/scripts
COPY --from=0 /app/magic-8ball /app/magic-8ball
COPY --from=0 /opt/pydeps /opt/pydeps

# Set environment so Python can find Playwright
ENV PYTHONPATH=/opt/pydeps:$PYTHONPATH
ENV PATH=/opt/pydeps/bin:$PATH

WORKDIR /app
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./magic-8ball"]
