FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o url-shortener ./cmd/http-server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/url-shortener .
CMD ["./url-shortener"]
