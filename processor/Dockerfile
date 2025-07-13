FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./src/*.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 4000

CMD ["./app"]