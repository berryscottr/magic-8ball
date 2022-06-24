FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

WORKDIR /app/pkg/bot
COPY * ./

WORKDIR /app/pkg/config
COPY *.go ./

WORKDIR /app

RUN go build main.go

EXPOSE 8080

CMD [ "./main" ]