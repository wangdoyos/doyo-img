# Build stage - Frontend
FROM node:20-alpine AS frontend
WORKDIR /app/web
COPY web/package.json web/package-lock.json* ./
RUN npm install --frozen-lockfile 2>/dev/null || npm install
COPY web/ .
RUN npm run build

# Build stage - Backend
FROM golang:1.24-alpine AS backend
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o doyo-img .

# Runtime stage
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend /app/doyo-img .
COPY config.example.yaml ./config.yaml

EXPOSE 9090
VOLUME ["/app/data"]

ENTRYPOINT ["./doyo-img"]
