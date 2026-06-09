package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/maz/k8s-mini-explorer/internal/services"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestEventService_List(t *testing.T) {
	client := fake.NewSimpleClientset(&corev1.Event{
		ObjectMeta: metav1.ObjectMeta{Name: "ev1", Namespace: "default"},
		Type:       "Warning",
		Reason:     "FailedScheduling",
		Message:    "0/1 nodes available",
		InvolvedObject: corev1.ObjectReference{
			Kind: "Pod",
			Name: "test-pod",
		},
		LastTimestamp: metav1.NewTime(time.Now()),
		Count:         2,
	})

	svc := services.NewEventService(client)
	resp, err := svc.List(context.Background(), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(resp.Events))
	}
	if resp.Events[0].Reason != "FailedScheduling" {
		t.Fatalf("unexpected reason: %s", resp.Events[0].Reason)
	}
}
