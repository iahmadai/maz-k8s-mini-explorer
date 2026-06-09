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
