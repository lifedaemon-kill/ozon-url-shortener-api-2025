FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/grpc-server ./cmd/grpc
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/http-gateway ./cmd/http-gw
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/swagger ./cmd/swagger


# grpc-api
FROM alpine:latest AS grpc
WORKDIR /app

COPY --from=builder /out/grpc-server    grpc-server
COPY --from=builder /app/migrations     migrations
COPY --from=builder /app/config         config
COPY --from=builder /app/.env           .env

ENTRYPOINT ["./grpc-server"]

# gateway-http-api
FROM alpine:latest AS gateway
WORKDIR /app

COPY --from=builder /out/http-gateway   http-gateway
COPY --from=builder /app/config         config
COPY --from=builder /app/.env           .env

ENTRYPOINT ["./http-gateway"]

# swagger
FROM alpine:latest AS swagger
WORKDIR /app

COPY --from=builder /out/swagger    swagger
COPY --from=builder /app/.env       .env
COPY --from=builder /app/config     config
COPY --from=builder /app/pkg  pkg

ENTRYPOINT ["./swagger"]
