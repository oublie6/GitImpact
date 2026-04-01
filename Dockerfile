FROM golang:1.22 AS builder

WORKDIR /src/backend
COPY backend/ ./

ENV CGO_ENABLED=0 \
    GOFLAGS=-mod=vendor

RUN go build -trimpath -ldflags "-s -w" -o /out/gitimpact-backend ./cmd/server

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /out/gitimpact-backend /app/gitimpact-backend

EXPOSE 8080
CMD ["/app/gitimpact-backend"]
