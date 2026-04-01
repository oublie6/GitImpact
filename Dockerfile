FROM golang:1.22

WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
COPY backend/vendor ./vendor
COPY backend/ ./

ENV CGO_ENABLED=0
ENV GOFLAGS=-mod=vendor

RUN go build -o /app/gitimpact-backend ./cmd/server

EXPOSE 8080
WORKDIR /app
CMD ["./gitimpact-backend"]
