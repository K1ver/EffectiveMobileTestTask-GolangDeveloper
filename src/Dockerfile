FROM golang:1.26.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/app ./cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/app /app
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/.env /.env
EXPOSE 8080
CMD ["/app"]