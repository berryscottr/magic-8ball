FROM golang:1.23-alpine as builder

ARG BOT_TOKEN
ENV BOT_TOKEN=$BOT_TOKEN

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN mkdir -p /app/data
WORKDIR /app/data
COPY data/* ./

RUN mkdir -p /app/pkg/bot
WORKDIR /app/pkg/bot
COPY pkg/bot/*.go ./

WORKDIR /app

RUN go build main.go

FROM alpine:latest as runtime

WORKDIR /app

COPY --from=builder /app/main /app/main

EXPOSE 8080

# HEALTHCHECK --interval=30s --timeout=10s \
#   CMD curl --fail http://localhost:8080/healthz || exit 1

CMD [ "./main" ]
