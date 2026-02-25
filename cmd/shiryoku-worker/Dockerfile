# Source: https://g3rhard.cc/posts/2025-12-01-docker-go/
# Caching go layer to make it more efficient

# ---- Build stage ----
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy dependency files first for layer caching
COPY go.mod go.sum ./

# Persist Go module cache between builds
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

ARG TARGETOS
ARG TARGETARCH

# Persist Go build cache + module cache
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" \
    -o shiryoku-worker ./cmd/shiryoku-worker


# ---- Runtime stage ----
FROM alpine:3.19

WORKDIR /app

RUN addgroup -S shiryoku && adduser -S shiryoku -G shiryoku

COPY --from=builder /app/shiryoku-worker .

USER shiryoku

ENTRYPOINT ["./shiryoku-worker"]
