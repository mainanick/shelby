FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build-linux

FROM debian:stretch-slim

WORKDIR /

COPY --from=builder /app/bin/unix/shellby /shellby

CMD ["/shellby"]
