package services_test

import (
	"context"
	"testing"

	"github.com/maz/k8s-mini-explorer/internal/services"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestServiceHealthService_Get(t *testing.T) {
	client := fake.NewSimpleClientset(
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{Name: "nginx", Namespace: "default"},
		},
		&corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{Name: "nginx", Namespace: "default"},
			Subsets: []corev1.EndpointSubset{
				{Addresses: []corev1.EndpointAddress{{IP: "10.0.0.5"}}},
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "nginx-pod", Namespace: "default"},
			Status: corev1.PodStatus{
				PodIP:  "10.0.0.5",
				PodIPs: []corev1.PodIP{{IP: "10.0.0.5"}},
				Phase:  corev1.PodRunning,
				Conditions: []corev1.PodCondition{
					{Type: corev1.PodReady, Status: corev1.ConditionTrue},
				},
			},
		},
	)

	svc := services.NewServiceHealthService(client)
	resp, err := svc.Get(context.Background(), "default", "nginx")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.EndpointCount != 1 {
		t.Fatalf("expected 1 endpoint, got %d", resp.EndpointCount)
	}
	if !resp.Healthy {
		t.Fatal("expected healthy service")
	}
	if resp.Endpoints[0].PodName != "nginx-pod" {
		t.Fatalf("unexpected pod name: %s", resp.Endpoints[0].PodName)
	}
}

func TestServiceHealthService_NoEndpoints(t *testing.T) {
	client := fake.NewSimpleClientset(
		&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "nginx", Namespace: "default"}},
		&corev1.Endpoints{ObjectMeta: metav1.ObjectMeta{Name: "nginx", Namespace: "default"}},
	)

	svc := services.NewServiceHealthService(client)
	resp, err := svc.Get(context.Background(), "default", "nginx")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Healthy {
		t.Fatal("expected unhealthy when no endpoints")
	}
}
