FROM golang:1.26.5-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata git
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG API_VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${APP_VERSION}" \
    -o /out/api ./cmd/api



FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /out/api /api

EXPOSE 8080
USER 65534:65534
ENTRYPOINT ["/api"]
