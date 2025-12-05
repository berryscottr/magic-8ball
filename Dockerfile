FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o magic-8ball ./main.go

FROM python:3.13-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    git wget curl ca-certificates --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

RUN pip install playwright
RUN playwright install --with-deps

COPY --from=0 /app/magic-8ball /app/magic-8ball
COPY --from=0 /app/data /app/data
COPY --from=0 /app/scripts /app/scripts

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./magic-8ball"]
