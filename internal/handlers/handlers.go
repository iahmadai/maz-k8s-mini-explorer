package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/maz/k8s-mini-explorer/internal/models"
	"github.com/maz/k8s-mini-explorer/internal/services"
	"github.com/maz/k8s-mini-explorer/pkg/apperrors"
)

type Handler struct {
	pods        *services.PodService
	deployments *services.DeploymentService
	events      *services.EventService
	services    *services.ServiceHealthService
}

func New(
	pods *services.PodService,
	deployments *services.DeploymentService,
	events *services.EventService,
	svcHealth *services.ServiceHealthService,
) *Handler {
	return &Handler{
		pods:        pods,
		deployments: deployments,
		events:      events,
		services:    svcHealth,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListPods(w http.ResponseWriter, r *http.Request) {
	ns := r.PathValue("namespace")
	resp, err := h.pods.List(r.Context(), ns)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetDeployment(w http.ResponseWriter, r *http.Request) {
	ns := r.PathValue("namespace")
	name := r.PathValue("name")
	resp, err := h.deployments.Get(r.Context(), ns, name)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ns := r.PathValue("namespace")
	resp, err := h.events.List(r.Context(), ns)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetServiceHealth(w http.ResponseWriter, r *http.Request) {
	ns := r.PathValue("namespace")
	name := r.PathValue("name")
	resp, err := h.services.Get(r.Context(), ns, name)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err error) {
	status := apperrors.HTTPStatus(err)
	writeJSON(w, status, models.ErrorResponse{Error: apperrors.Message(err)})
}
