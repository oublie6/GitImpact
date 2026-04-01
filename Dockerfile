FROM node:20-bookworm-slim AS frontend-builder

WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build:offline

FROM golang:1.22 AS backend-builder

WORKDIR /src/backend
COPY backend/ ./

ENV CGO_ENABLED=0 \
    GOFLAGS=-mod=vendor

RUN go build -trimpath -ldflags "-s -w" -o /out/gitimpact-backend ./cmd/server

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=backend-builder /out/gitimpact-backend /app/gitimpact-backend
COPY --from=backend-builder /src/backend/config.example.yaml /app/config.example.yaml
COPY --from=frontend-builder /src/frontend/dist /app/web/dist

EXPOSE 8080
CMD ["/app/gitimpact-backend"]
