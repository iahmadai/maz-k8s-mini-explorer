package models

type PodSummary struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Age    string `json:"age"`
	Ready  string `json:"ready"`
}

type PodListResponse struct {
	Namespace string       `json:"namespace"`
	Pods      []PodSummary `json:"pods"`
}

type DeploymentCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

type DeploymentStatusResponse struct {
	Name                  string                `json:"name"`
	Namespace             string                `json:"namespace"`
	Replicas              int32                 `json:"replicas"`
	ReadyReplicas         int32                 `json:"ready_replicas"`
	AvailableReplicas     int32                 `json:"available_replicas"`
	UnavailableReplicas   int32                 `json:"unavailable_replicas"`
	UpdatedReplicas       int32                 `json:"updated_replicas"`
	Conditions            []DeploymentCondition `json:"conditions"`
}

type EventSummary struct {
	Type           string `json:"type"`
	Reason         string `json:"reason"`
	Message        string `json:"message"`
	InvolvedObject string `json:"involved_object"`
	Age            string `json:"age"`
	Count          int32  `json:"count"`
}

type EventListResponse struct {
	Namespace string         `json:"namespace"`
	Events    []EventSummary `json:"events"`
}

type EndpointHealth struct {
	IP        string `json:"ip"`
	PodName   string `json:"pod_name,omitempty"`
	PodStatus string `json:"pod_status,omitempty"`
	Ready     bool   `json:"ready"`
}

type ServiceHealthResponse struct {
	Name          string           `json:"name"`
	Namespace     string           `json:"namespace"`
	EndpointCount int              `json:"endpoint_count"`
	Endpoints     []EndpointHealth `json:"endpoints"`
	Healthy       bool             `json:"healthy"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
