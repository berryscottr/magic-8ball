FROM golang:1.21-alpine

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

EXPOSE 8080

CMD [ "./main" ]
