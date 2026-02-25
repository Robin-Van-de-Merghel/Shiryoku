# Source: https://g3rhard.cc/posts/2025-12-01-docker-go/
# Caching go layer to make it more efficient

# ---- Build stage ----
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /app

# Install git only if your modules require it
RUN apk add --no-cache git

# Copy go mod files first (better layer caching)
COPY go.mod go.sum ./

# Cache Go modules between builds
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the rest of the source
COPY . .

ARG TARGETOS
ARG TARGETARCH

# Cache Go build cache between builds
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" \
    -o shiryoku-api ./cmd/shiryoku-api


# ---- Runtime stage ----
FROM alpine:3.19

WORKDIR /app

RUN addgroup -S shiryoku && adduser -S shiryoku -G shiryoku
RUN apk add --no-cache curl

COPY --from=builder /app/shiryoku-api .

USER shiryoku

EXPOSE 8080

ENTRYPOINT ["./shiryoku-api"]
