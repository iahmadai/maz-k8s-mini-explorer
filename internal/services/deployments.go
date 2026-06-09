package services

import (
	"context"

	"github.com/maz/k8s-mini-explorer/internal/models"
	"github.com/maz/k8s-mini-explorer/pkg/apperrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentService struct {
	client kubernetes.Interface
}

func NewDeploymentService(client kubernetes.Interface) *DeploymentService {
	return &DeploymentService{client: client}
}

func (s *DeploymentService) Get(ctx context.Context, namespace, name string) (*models.DeploymentStatusResponse, error) {
	if namespace == "" || name == "" {
		return nil, apperrors.BadRequest("namespace and deployment name are required")
	}

	dep, err := s.client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, apperrors.FromK8s(err)
	}

	replicas := int32(0)
	if dep.Spec.Replicas != nil {
		replicas = *dep.Spec.Replicas
	}

	conditions := make([]models.DeploymentCondition, 0, len(dep.Status.Conditions))
	for _, c := range dep.Status.Conditions {
		conditions = append(conditions, models.DeploymentCondition{
			Type:    string(c.Type),
			Status:  string(c.Status),
			Reason:  c.Reason,
			Message: c.Message,
		})
	}

	return &models.DeploymentStatusResponse{
		Name:                dep.Name,
		Namespace:           dep.Namespace,
		Replicas:            replicas,
		ReadyReplicas:       dep.Status.ReadyReplicas,
		AvailableReplicas:   dep.Status.AvailableReplicas,
		UnavailableReplicas: dep.Status.UnavailableReplicas,
		UpdatedReplicas:     dep.Status.UpdatedReplicas,
		Conditions:          conditions,
	}, nil
}
