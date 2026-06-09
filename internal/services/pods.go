package services

import (
	"context"
	"fmt"

	"github.com/maz/k8s-mini-explorer/internal/models"
	"github.com/maz/k8s-mini-explorer/pkg/apperrors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodService struct {
	client kubernetes.Interface
}

func NewPodService(client kubernetes.Interface) *PodService {
	return &PodService{client: client}
}

func (s *PodService) List(ctx context.Context, namespace string) (*models.PodListResponse, error) {
	if namespace == "" {
		return nil, apperrors.BadRequest("namespace is required")
	}

	list, err := s.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, apperrors.FromK8s(err)
	}

	pods := make([]models.PodSummary, 0, len(list.Items))
	for _, pod := range list.Items {
		pods = append(pods, models.PodSummary{
			Name:   pod.Name,
			Status: podStatus(&pod),
			Age:    formatAge(pod.CreationTimestamp.Time),
			Ready:  podReady(&pod),
		})
	}

	return &models.PodListResponse{
		Namespace: namespace,
		Pods:      pods,
	}, nil
}

func podStatus(pod *corev1.Pod) string {
	if pod.Status.Phase != "" {
		return string(pod.Status.Phase)
	}
	return "Unknown"
}

func podReady(pod *corev1.Pod) string {
	ready := 0
	total := len(pod.Status.ContainerStatuses)
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Ready {
			ready++
		}
	}
	if total == 0 {
		return fmt.Sprintf("0/%d", len(pod.Spec.Containers))
	}
	return fmt.Sprintf("%d/%d", ready, total)
}
