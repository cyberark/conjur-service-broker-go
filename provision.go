package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServiceInstanceDeprovision deprovision a service instance
// (DELETE /v2/service_instances/{instance_id})
func (*ServerImpl) ServiceInstanceDeprovision(c *gin.Context, instanceID string, params ServiceInstanceDeprovisionParams) {
	// That's all folks!
	c.Status(http.StatusOK)
}

// ServiceInstanceGet get a service instance
// (GET /v2/service_instances/{instance_id})
func (*ServerImpl) ServiceInstanceGet(c *gin.Context, instanceID string, params ServiceInstanceGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}

// ServiceInstanceUpdate update a service instance
// (PATCH /v2/service_instances/{instance_id})
func (*ServerImpl) ServiceInstanceUpdate(c *gin.Context, instanceID string, params ServiceInstanceUpdateParams) {
	// That's all folks!
	c.Status(http.StatusOK)
}

// ServiceInstanceProvision provision a service instance
// (PUT /v2/service_instances/{instance_id})
func (s *ServerImpl) ServiceInstanceProvision(c *gin.Context, instanceID string, params ServiceInstanceProvisionParams) {
	body := ServiceInstanceProvisionJSONRequestBody{}
	err := c.BindJSON(&body)
	if err != nil {
		// TODO: handle error from AbortWithError
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse request body: %w", err))
	}
	if err = validateServiceAndPlan(body.ServiceId, body.PlanId); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("validation failed: %w", err))
	}

	// check if exists
	// TODO: use IDs from context
	orgSpace := s.client.NewOrgSpace(
		body.OrganizationGuid,
		body.SpaceGuid,
		formContext(body.Context, "organization_name"),
		formContext(body.Context, "space_name"),
	)
	if err = orgSpace.CreatePolicy(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create policy: %w", err))
	}
	if err = orgSpace.ValidateExists(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to validate policy exists: %w", err))
	}

	c.Status(http.StatusCreated)
}

// ServiceInstanceLastOperationGet get the last requested operation state for service instance
// (GET /v2/service_instances/{instance_id}/last_operation)
func (*ServerImpl) ServiceInstanceLastOperationGet(c *gin.Context, instanceID string, params ServiceInstanceLastOperationGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}
