FROM golang:1.18-alpine

WORKDIR /app

ADD go.mod ./
ADD go.sum ./
RUN go mod download

ADD *.go ./

WORKDIR /app/pkg/bot
ADD * ./

WORKDIR /app/pkg/config
ADD *.go ./

WORKDIR /app

RUN go build main.go

EXPOSE 8080

CMD [ "./main" ]