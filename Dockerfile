FROM golang:1.23-alpine

ARG BOT_TOKEN
ENV BOT_TOKEN=$BOT_TOKEN

RUN apk add --no-cache git ca-certificates wget

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o magic-8ball ./main.go

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./magic-8ball"]
