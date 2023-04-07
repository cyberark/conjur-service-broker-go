package servicebroker

import (
	"fmt"
	"net/http"

	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/gin-gonic/gin"
)

// ServiceInstanceDeprovision deprovision a service instance
// (DELETE /v2/service_instances/{instance_id})
func (*server) ServiceInstanceDeprovision(c *gin.Context, instanceID string, params ServiceInstanceDeprovisionParams) {
	if err := validateServiceAndPlan(params.ServiceId, &params.PlanId); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("validation failed: %w", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// ServiceInstanceGet get a service instance
// (GET /v2/service_instances/{instance_id})
func (*server) ServiceInstanceGet(c *gin.Context, instanceID string, params ServiceInstanceGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}

// ServiceInstanceUpdate update a service instance
// (PATCH /v2/service_instances/{instance_id})
func (*server) ServiceInstanceUpdate(c *gin.Context, instanceID string, params ServiceInstanceUpdateParams) {
	body := ServiceInstanceUpdateJSONRequestBody{}
	err := c.BindJSON(&body)
	if err != nil {
		// TODO: handle error from AbortWithError
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse request body: %w", err))
		return
	}
	if err = validateServiceAndPlan(body.ServiceId, body.PlanId); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("validation failed: %w", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// ServiceInstanceProvision provision a service instance
// (PUT /v2/service_instances/{instance_id})
func (s *server) ServiceInstanceProvision(c *gin.Context, instanceID string, params ServiceInstanceProvisionParams) {
	body := ServiceInstanceProvisionJSONRequestBody{}
	err := c.BindJSON(&body)
	if err != nil {
		// TODO: handle error from AbortWithError
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse request body: %w", err))
		return
	}
	if err = validateServiceAndPlan(body.ServiceId, &body.PlanId); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("validation failed: %w", err))
		return
	}

	ctxParams := parseContext(body.Context)
	orgSpace := conjur.NewProvision(
		s.client,
		ctxParams.OrgID,
		ctxParams.SpaceID,
		ctxParams.OrgName,
		ctxParams.SpaceName,
		s.enableSpaceIdentity,
	)
	if err = orgSpace.CreatePolicy(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create policy: %w", err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// ServiceInstanceLastOperationGet get the last requested operation state for service instance
// (GET /v2/service_instances/{instance_id}/last_operation)
func (*server) ServiceInstanceLastOperationGet(c *gin.Context, instanceID string, params ServiceInstanceLastOperationGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}
