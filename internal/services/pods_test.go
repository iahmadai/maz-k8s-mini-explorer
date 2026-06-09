package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/maz/k8s-mini-explorer/internal/services"
	"github.com/maz/k8s-mini-explorer/pkg/apperrors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestPodService_List(t *testing.T) {
	client := fake.NewSimpleClientset(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "nginx-abc",
			Namespace:         "default",
			CreationTimestamp: metav1.NewTime(time.Now().Add(-2 * time.Hour)),
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			ContainerStatuses: []corev1.ContainerStatus{
				{Ready: true},
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{Name: "nginx"}},
		},
	})

	svc := services.NewPodService(client)
	resp, err := svc.List(context.Background(), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Pods) != 1 {
		t.Fatalf("expected 1 pod, got %d", len(resp.Pods))
	}
	if resp.Pods[0].Name != "nginx-abc" {
		t.Fatalf("unexpected pod name: %s", resp.Pods[0].Name)
	}
	if resp.Pods[0].Status != "Running" {
		t.Fatalf("unexpected status: %s", resp.Pods[0].Status)
	}
}

func TestPodService_ListEmptyNamespace(t *testing.T) {
	svc := services.NewPodService(fake.NewSimpleClientset())
	_, err := svc.List(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty namespace")
	}
	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) || appErr.Code != 400 {
		t.Fatalf("expected 400 bad request, got %v", err)
	}
}
