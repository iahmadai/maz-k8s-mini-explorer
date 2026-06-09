package services

import (
	"context"
	"sort"
	"time"

	"github.com/maz/k8s-mini-explorer/internal/models"
	"github.com/maz/k8s-mini-explorer/pkg/apperrors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type EventService struct {
	client kubernetes.Interface
}

func NewEventService(client kubernetes.Interface) *EventService {
	return &EventService{client: client}
}

func (s *EventService) List(ctx context.Context, namespace string) (*models.EventListResponse, error) {
	if namespace == "" {
		return nil, apperrors.BadRequest("namespace is required")
	}

	list, err := s.client.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, apperrors.FromK8s(err)
	}

	items := list.Items
	sort.Slice(items, func(i, j int) bool {
		return eventTime(items[i]).After(eventTime(items[j]))
	})

	events := make([]models.EventSummary, 0, len(items))
	for _, ev := range items {
		events = append(events, models.EventSummary{
			Type:           ev.Type,
			Reason:         ev.Reason,
			Message:        ev.Message,
			InvolvedObject: ev.InvolvedObject.Kind + "/" + ev.InvolvedObject.Name,
			Age:            formatAge(eventTime(ev)),
			Count:          ev.Count,
		})
	}

	return &models.EventListResponse{
		Namespace: namespace,
		Events:    events,
	}, nil
}

func eventTime(ev corev1.Event) time.Time {
	if !ev.LastTimestamp.IsZero() {
		return ev.LastTimestamp.Time
	}
	if !ev.EventTime.IsZero() {
		return ev.EventTime.Time
	}
	return ev.CreationTimestamp.Time
}
