# ---- Build stage ----
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" \
    -o shiryoku-worker ./cmd/shiryoku-worker

# ---- Runtime stage ----
FROM alpine:3.19

WORKDIR /app

RUN addgroup -S shiryoku && adduser -S shiryoku -G shiryoku

COPY --from=builder /app/shiryoku-worker .

USER shiryoku

ENTRYPOINT ["./shiryoku-worker"]
