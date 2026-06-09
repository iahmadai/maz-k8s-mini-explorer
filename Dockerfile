# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /k8s-explorer ./cmd/k8s-explorer

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /k8s-explorer /k8s-explorer

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/k8s-explorer"]
CMD ["serve"]

# Dev compose stage (kubeconfig fix for Docker Desktop)
FROM alpine:3.20 AS dev
RUN apk add --no-cache ca-certificates
COPY --from=builder /k8s-explorer /k8s-explorer
COPY scripts/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
CMD ["serve"]
