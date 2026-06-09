package server

import (
	"log/slog"
	"net/http"

	"github.com/maz/k8s-mini-explorer/internal/handlers"
	"github.com/maz/k8s-mini-explorer/internal/services"
	"k8s.io/client-go/kubernetes"
)

func NewRouter(client kubernetes.Interface) http.Handler {
	h := handlers.New(
		services.NewPodService(client),
		services.NewDeploymentService(client),
		services.NewEventService(client),
		services.NewServiceHealthService(client),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /api/v1/namespaces/{namespace}/pods", h.ListPods)
	mux.HandleFunc("GET /api/v1/namespaces/{namespace}/deployments/{name}", h.GetDeployment)
	mux.HandleFunc("GET /api/v1/namespaces/{namespace}/events", h.ListEvents)
	mux.HandleFunc("GET /api/v1/namespaces/{namespace}/services/{name}/health", h.GetServiceHealth)

	return loggingMiddleware(mux)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
