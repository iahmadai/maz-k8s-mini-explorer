# Architecture

## Overview

`maz-k8s-mini-explorer` is a dual-entry Go application: a Cobra CLI and an HTTP REST API share the same service layer and Kubernetes client abstraction.

```
cmd/k8s-explorer/main.go
        │
        ├── internal/cli/          (cobra commands)
        └── internal/server/       (HTTP router)
                │
                └── internal/handlers/
                        │
                        └── internal/services/   (business logic)
                                │
                                └── internal/k8s/  (client-go wrapper)
                                        │
                                        └── Kubernetes API
```

## Design Decisions

### Shared service layer

All Kubernetes logic lives in `internal/services/`. Both the CLI and HTTP handlers call the same services, ensuring consistent behavior and a single place to test.

### Kubernetes client interface

Services depend on `kubernetes.Interface` from client-go, not a concrete `*kubernetes.Clientset`. This allows unit tests to inject `fake.NewSimpleClientset()` without a real cluster.

```go
// internal/k8s/client.go
func NewClient(cfg *config.Config) (kubernetes.Interface, error)
```

Connection modes:
- **Local dev**: kubeconfig file (`KUBECONFIG` or `~/.kube/config`)
- **In-cluster**: ServiceAccount token when `IN_CLUSTER=true`

### Single binary

One entry point (`cmd/k8s-explorer`) exposes all CLI subcommands plus `serve`. This avoids duplicating initialization code across separate binaries.

### HTTP routing

Uses Go 1.22+ stdlib `http.ServeMux` with method-aware path patterns instead of a third-party router. Middleware handles request logging.

### Error mapping

`pkg/apperrors` translates errors to HTTP status codes:
- Kubernetes `IsNotFound` → 404
- Kubernetes `IsForbidden` → 403
- Validation errors → 400
- Other API errors → 500

Services return `*apperrors.AppError` for domain validation; Kubernetes errors are wrapped via `apperrors.FromK8s()`.

### Service health resolution

Service health checks:
1. Fetch `Endpoints` for the service name
2. List pods in the namespace
3. Match endpoint IPs to pod IPs
4. Mark healthy when all endpoints are ready and count > 0

### Configuration

`internal/config/config.go` loads from environment variables with sensible defaults. CLI persistent flags override at runtime.

| Source | Priority |
|--------|----------|
| CLI flags | Highest |
| Environment variables | Medium |
| Defaults | Lowest |

## Project Structure

```
maz-k8s-mini-explorer/
├── cmd/k8s-explorer/       # main entry
├── internal/
│   ├── cli/                # cobra commands + output formatting
│   ├── config/             # env/flag configuration
│   ├── handlers/           # HTTP handlers
│   ├── k8s/                # client-go factory
│   ├── models/             # response DTOs
│   ├── server/             # router + middleware
│   └── services/           # core logic + tests
├── pkg/apperrors/          # typed errors → HTTP status
├── deploy/rbac.yaml        # in-cluster deployment + RBAC
├── docs/                   # API and architecture docs
├── Dockerfile              # multi-stage build (distroless)
└── docker-compose.yml      # local dev with kubeconfig mount
```

## Testing Strategy

| Layer | Approach |
|-------|----------|
| Services | `fake.NewSimpleClientset()` with pre-seeded resources |
| Handlers | `httptest` with fake client injected via services |
| CLI | Manual / integration (services covered by unit tests) |

Each service test covers: happy path, not found, and empty list where applicable.

## Dependencies

| Package | Purpose |
|---------|---------|
| `k8s.io/client-go` | Kubernetes API client |
| `k8s.io/api` | Kubernetes resource types |
| `github.com/spf13/cobra` | CLI framework |

HTTP uses only the Go standard library (`net/http`, `encoding/json`).

## Deployment

For in-cluster use:
- Set `IN_CLUSTER=true`
- Apply `deploy/rbac.yaml` (read-only ClusterRole)
- Image runs as non-root (distroless)

For local Docker:
- Mount host kubeconfig read-only
- Set `KUBECONFIG=/root/.kube/config`

## Future Considerations

- Context timeout on Kubernetes API calls
- Pagination for large pod/event lists
- OpenAPI/Swagger spec generation
- CLI integration tests with captured stdout
