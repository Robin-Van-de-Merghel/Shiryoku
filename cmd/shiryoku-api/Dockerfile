# ---- Build stage ----
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (required for some Go modules)
RUN apk add --no-cache git

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build arguments for cross-compilation
ARG TARGETOS
ARG TARGETARCH

# Build static binary
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" \
    -o shiryoku-api ./cmd/shiryoku-api

# ---- Runtime stage ----
FROM alpine:3.19

WORKDIR /app

# Add non-root user
RUN addgroup -S shiryoku && adduser -S shiryoku -G shiryoku
RUN apk add --no-cache curl

# Copy binary
COPY --from=builder /app/shiryoku-api .

USER shiryoku

EXPOSE 8080

ENTRYPOINT ["./shiryoku-api"]
