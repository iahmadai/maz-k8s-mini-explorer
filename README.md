# maz-k8s-mini-explorer

A lightweight Go tool for exploring Kubernetes clusters. Provides both a **REST API** and a **CLI** to inspect pods, deployments, events, and service health.

## Features

- List pods in a namespace (name, status, age, ready count)
- Get deployment status (replicas, conditions)
- List namespace events (sorted by last timestamp)
- Check service health via endpoints and pod mapping
- REST API (`/api/v1/...`) and Cobra CLI with JSON or table output
- Unit tests with `fake.Clientset` (no live cluster required)
- Docker image and in-cluster deployment manifests

## Prerequisites

- Go 1.23+
- Access to a Kubernetes cluster (local or remote)
- `kubectl` configured (for local dev) or in-cluster ServiceAccount (for deployment)

### Local cluster setup (optional)

If you don't have a cluster, use [kind](https://kind.sigs.k8s.io/) or [minikube](https://minikube.sigs.k8s.io/):

```bash
# kind
kind create cluster

# minikube
minikube start
```

## Quick Start

```bash
# Clone and build
make build

# CLI — list pods
./bin/k8s-explorer pods --namespace default --output table

# CLI — deployment status
./bin/k8s-explorer deployment nginx --namespace default

# CLI — events
./bin/k8s-explorer events --namespace default --output table

# CLI — service health
./bin/k8s-explorer service-health nginx --namespace default

# Start REST API server
./bin/k8s-explorer serve --port 8080
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `KUBECONFIG` | `~/.kube/config` | Path to kubeconfig file |
| `K8S_CONTEXT` | (current context) | Kubernetes context name |
| `K8S_NAMESPACE` | `default` | Default namespace for CLI |
| `PORT` | `8080` | HTTP port for `serve` |
| `IN_CLUSTER` | `false` | Use in-cluster config (set `true` inside a pod) |

CLI flags override env vars: `--kubeconfig`, `--context`, `--namespace`, `--in-cluster`, `--output`.

## Docker

Run against your local kubeconfig:

```bash
docker compose up --build
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/namespaces/default/pods
```

The compose file mounts `~/.kube` read-only into the container.

## Deploy to Kubernetes

Build and load the image (example with kind):

```bash
docker build -t k8s-mini-explorer:latest .
kind load docker-image k8s-mini-explorer:latest
kubectl apply -f deploy/rbac.yaml
kubectl port-forward svc/k8s-mini-explorer 8080:8080
```

The manifest includes a read-only ClusterRole (pods, deployments, events, endpoints, services), ServiceAccount, Deployment, and Service.

## Development

```bash
make test      # run all tests with coverage
make build     # build bin/k8s-explorer
make serve     # build and start API on :8080
make tidy      # go mod tidy
```

## API Documentation

See [docs/API.md](docs/API.md) for endpoint details and sample responses.

## Architecture

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for design decisions and project structure.

## License

MIT
