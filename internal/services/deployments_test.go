package services_test

import (
	"context"
	"testing"

	"github.com/maz/k8s-mini-explorer/internal/services"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDeploymentService_Get(t *testing.T) {
	replicas := int32(3)
	client := fake.NewSimpleClientset(&appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "nginx", Namespace: "default"},
		Spec:       appsv1.DeploymentSpec{Replicas: &replicas},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas:       3,
			AvailableReplicas:   3,
			UnavailableReplicas: 0,
			UpdatedReplicas:     3,
		},
	})

	svc := services.NewDeploymentService(client)
	resp, err := svc.Get(context.Background(), "default", "nginx")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Replicas != 3 || resp.ReadyReplicas != 3 {
		t.Fatalf("unexpected replicas: %+v", resp)
	}
}

func TestDeploymentService_NotFound(t *testing.T) {
	svc := services.NewDeploymentService(fake.NewSimpleClientset())
	_, err := svc.Get(context.Background(), "default", "missing")
	if err == nil {
		t.Fatal("expected not found error")
	}
}
