FROM golang:1.24.1-alpine as builder

RUN mkdir /app

ADD . /app

WORKDIR /app

COPY go.mod go.sum .env dict.json ./

RUN go mod tidy

COPY . .

RUN go build ./main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 2112

CMD ["./main"]

