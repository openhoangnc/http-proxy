FROM golang:alpine AS builder

WORKDIR /go/src/app

COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go build -ldflags="-s -w" -o /app .

FROM scratch

COPY --from=builder /app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app"]
