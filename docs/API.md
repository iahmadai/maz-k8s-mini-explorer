# API Reference

Base URL: `http://localhost:8080`

All responses are JSON. Errors return `{"error": "message"}` with an appropriate HTTP status code.

## Health

### `GET /health`

Liveness check.

**Response `200`**

```json
{
  "status": "ok"
}
```

---

## Pods

### `GET /api/v1/namespaces/{namespace}/pods`

List all pods in the given namespace.

**Response `200`**

```json
{
  "namespace": "default",
  "pods": [
    {
      "name": "nginx-abc123",
      "status": "Running",
      "age": "2h15m",
      "ready": "1/1"
    }
  ]
}
```

| Field | Description |
|-------|-------------|
| `status` | Pod phase (`Running`, `Pending`, etc.) |
| `age` | Human-readable time since creation |
| `ready` | Ready containers / total containers |

---

## Deployments

### `GET /api/v1/namespaces/{namespace}/deployments/{name}`

Get deployment status.

**Response `200`**

```json
{
  "name": "nginx",
  "namespace": "default",
  "replicas": 3,
  "ready_replicas": 3,
  "available_replicas": 3,
  "unavailable_replicas": 0,
  "updated_replicas": 3,
  "conditions": [
    {
      "type": "Available",
      "status": "True",
      "reason": "MinimumReplicasAvailable",
      "message": "Deployment has minimum availability."
    }
  ]
}
```

---

## Events

### `GET /api/v1/namespaces/{namespace}/events`

List events in the namespace, sorted by last timestamp (newest first).

**Response `200`**

```json
{
  "namespace": "default",
  "events": [
    {
      "type": "Normal",
      "reason": "Scheduled",
      "message": "Successfully assigned default/nginx to node-1",
      "involved_object": "Pod/nginx-abc123",
      "age": "5m",
      "count": 1
    }
  ]
}
```

---

## Service Health

### `GET /api/v1/namespaces/{namespace}/services/{name}/health`

Check service health by resolving endpoints and matching pods by IP.

**Response `200`**

```json
{
  "name": "nginx",
  "namespace": "default",
  "endpoint_count": 3,
  "endpoints": [
    {
      "ip": "10.0.0.5",
      "pod_name": "nginx-xxx",
      "pod_status": "Running",
      "ready": true
    }
  ],
  "healthy": true
}
```

| Field | Description |
|-------|-------------|
| `healthy` | `true` when at least one endpoint exists and all are ready |
| `endpoints` | Resolved IPs with optional pod name/status |

---

## Error Responses

| Status | Scenario |
|--------|----------|
| `400` | Empty namespace or resource name |
| `403` | Forbidden by Kubernetes RBAC |
| `404` | Resource not found |
| `500` | Internal / unexpected Kubernetes API error |
| `503` | Cluster unreachable |

**Example `404`**

```json
{
  "error": "deployments.apps \"missing\" not found"
}
```

## CLI Equivalents

| API | CLI |
|-----|-----|
| `GET .../pods` | `k8s-explorer pods --namespace {ns}` |
| `GET .../deployments/{name}` | `k8s-explorer deployment {name} --namespace {ns}` |
| `GET .../events` | `k8s-explorer events --namespace {ns}` |
| `GET .../services/{name}/health` | `k8s-explorer service-health {name} --namespace {ns}` |
| (server) | `k8s-explorer serve --port 8080` |

Add `--output table` for human-readable CLI output instead of JSON.
