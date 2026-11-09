package handler

import (
	"net/http"

	"github.com/Baroneee/EBVN_Assignment_1/internal/service"
	"github.com/gin-gonic/gin"
)

// HealthCheck handles health check HTTP requests.
type HealthCheck interface {
	Check(c *gin.Context)
}

// healthCheck is the HTTP handler for health check requests.
type healthCheck struct {
	svc service.HealthCheck
}

// NewHealthCheck creates a health check handler backed by svc.
func NewHealthCheck(svc service.HealthCheck) HealthCheck {
	return &healthCheck{
		svc: svc,
	}
}

// @Summary Check the health of the service
// @Description Returns the health status of the service
// @Tags heatlh-check
// @Produce application/json
// @Success 200 {object} service.HealthResponse
// @Failure 500 {object} string "Internal Server Error"
// @Router /health [get]
// Check writes the service health status as a JSON response.
func (h *healthCheck) Check(c *gin.Context) {
	result := h.svc.Check()
	c.JSON(http.StatusOK, result)
}
