package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maz/k8s-mini-explorer/internal/services"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func newTestHandler() *Handler {
	client := fake.NewSimpleClientset(
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "default"},
			Status:     corev1.PodStatus{Phase: corev1.PodRunning},
		},
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: "default"},
		},
	)
	return New(
		services.NewPodService(client),
		services.NewDeploymentService(client),
		services.NewEventService(client),
		services.NewServiceHealthService(client),
	)
}

func TestHandler_Health(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.Health(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_ListPods(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/namespaces/default/pods", nil)
	req.SetPathValue("namespace", "default")
	rec := httptest.NewRecorder()
	h.ListPods(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestHandler_GetDeploymentNotFound(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/namespaces/default/deployments/missing", nil)
	req.SetPathValue("namespace", "default")
	req.SetPathValue("name", "missing")
	rec := httptest.NewRecorder()
	h.GetDeployment(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
