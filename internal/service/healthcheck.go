package service

import (
	"uuid"
)

// HealthResponse contains the health status returned by the service.
//
//go:generate mockery --name=HealthCheck --filename=healthcheck.go --outpkg svc_mocks
type HealthResponse struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}

// HealthCheck provides the application's health status.
type HealthCheck interface {
	Check() HealthResponse
}

// healthCheck stores the service name and instance identifier used in responses.
type healthCheck struct {
	serviceName string
	instanceID  string
}

// NewHealthCheck creates a health check service, generating an instance ID when needed.
func NewHealthCheck(serviceName, instanceID string) HealthCheck {
	if instanceID == "" {
		instanceID = uuid.New().String()
	}

	return &healthCheck{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

// Check returns the current health status.
func (h *healthCheck) Check() HealthResponse {
	return HealthResponse{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}
}
