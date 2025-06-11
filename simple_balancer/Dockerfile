FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o loadbalancer ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/loadbalancer .
COPY ./cmd/config.json .
EXPOSE 8080
CMD ["./loadbalancer", "--config", "config.json"]
