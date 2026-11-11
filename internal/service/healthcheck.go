package service

import (
	"context"
	"errors"
	"uuid"

	"github.com/redis/go-redis/v9"
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
	Check() (HealthResponse, error)
}

// healthCheck stores the service name and instance identifier used in responses.
type healthCheck struct {
	serviceName string
	instanceID  string
	redisClient *redis.Client
}

var errRedisClientUnavailable = errors.New("redis client unavailable")

// NewHealthCheck creates a health check service, generating an instance ID when needed.
func NewHealthCheck(serviceName, instanceID string, redisClient *redis.Client) HealthCheck {
	if instanceID == "" {
		instanceID = uuid.New().String()
	}

	return &healthCheck{
		serviceName: serviceName,
		instanceID:  instanceID,
		redisClient: redisClient,
	}
}

// Check returns the current health status.
func (h *healthCheck) Check() (HealthResponse, error) {
	if h.redisClient == nil {
		return HealthResponse{
			ServiceName: h.serviceName,
			InstanceID:  h.instanceID,
		}, errRedisClientUnavailable
	}

	if err := h.redisClient.Ping(context.Background()).Err(); err != nil {
		return HealthResponse{
			ServiceName: h.serviceName,
			InstanceID:  h.instanceID,
		}, err
	}

	return HealthResponse{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}, nil
}
