FROM golang:latest AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o monerochat .

EXPOSE 8080

CMD ["./monerochat"]

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/src/app/monerochat .

EXPOSE 5000

CMD ["./monerochat"]
