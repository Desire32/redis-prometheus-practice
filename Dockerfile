FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main .

FROM gcr.io/distroless/base-debian11

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

CMD ["./main"]
