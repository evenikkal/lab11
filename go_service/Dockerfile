# Task A1: static Go compilation for minimal scratch image.
# CGO_ENABLED=0 produces a fully self-contained binary with no libc dependency.

# ── build stage ────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go_service .

# ── runtime stage: scratch (zero-overhead base image) ──────────────────
FROM scratch
COPY --from=builder /go_service /go_service
EXPOSE 8082
ENTRYPOINT ["/go_service"]
