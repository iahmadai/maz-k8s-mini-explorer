package services

import (
	"context"

	"github.com/maz/k8s-mini-explorer/internal/models"
	"github.com/maz/k8s-mini-explorer/pkg/apperrors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceHealthService struct {
	client kubernetes.Interface
}

func NewServiceHealthService(client kubernetes.Interface) *ServiceHealthService {
	return &ServiceHealthService{client: client}
}

func (s *ServiceHealthService) Get(ctx context.Context, namespace, name string) (*models.ServiceHealthResponse, error) {
	if namespace == "" || name == "" {
		return nil, apperrors.BadRequest("namespace and service name are required")
	}

	if _, err := s.client.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{}); err != nil {
		return nil, apperrors.FromK8s(err)
	}

	endpoints, err := s.client.CoreV1().Endpoints(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, apperrors.FromK8s(err)
	}

	podList, err := s.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, apperrors.FromK8s(err)
	}

	ipToPod := map[string]*corev1.Pod{}
	for i := range podList.Items {
		pod := &podList.Items[i]
		for _, ip := range pod.Status.PodIPs {
			ipToPod[ip.IP] = pod
		}
		if pod.Status.PodIP != "" {
			ipToPod[pod.Status.PodIP] = pod
		}
	}

	result := make([]models.EndpointHealth, 0)
	allReady := true

	for _, subset := range endpoints.Subsets {
		for _, addr := range subset.Addresses {
			ep := models.EndpointHealth{IP: addr.IP, Ready: true}
			if pod, ok := ipToPod[addr.IP]; ok {
				ep.PodName = pod.Name
				ep.PodStatus = podStatus(pod)
				ep.Ready = isPodReady(pod)
			}
			if !ep.Ready {
				allReady = false
			}
			result = append(result, ep)
		}
		for _, addr := range subset.NotReadyAddresses {
			allReady = false
			ep := models.EndpointHealth{IP: addr.IP, Ready: false}
			if pod, ok := ipToPod[addr.IP]; ok {
				ep.PodName = pod.Name
				ep.PodStatus = podStatus(pod)
			}
			result = append(result, ep)
		}
	}

	healthy := len(result) > 0 && allReady

	return &models.ServiceHealthResponse{
		Name:          name,
		Namespace:     namespace,
		EndpointCount: len(result),
		Endpoints:     result,
		Healthy:       healthy,
	}, nil
}

func isPodReady(pod *corev1.Pod) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == corev1.PodReady {
			return cond.Status == corev1.ConditionTrue
		}
	}
	return false
}
