// Package servicebroker provides an implementation of the generated service broker server
package servicebroker

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.cyberng.com/Conjur-Enterprise/conjur-service-broker-go/internal/ctxutil"
)

// ServiceInstanceDeprovision deprovision a service instance
// (DELETE /v2/service_instances/{instance_id})
func (*server) ServiceInstanceDeprovision(c *gin.Context, _ string, _ ServiceInstanceDeprovisionParams) {
	c.JSON(http.StatusOK, gin.H{})
}

// ServiceInstanceGet get a service instance
// (GET /v2/service_instances/{instance_id})
func (*server) ServiceInstanceGet(c *gin.Context, _ string, _ ServiceInstanceGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}

// ServiceInstanceUpdate update a service instance
// (PATCH /v2/service_instances/{instance_id})
func (*server) ServiceInstanceUpdate(c *gin.Context, _ string, _ ServiceInstanceUpdateParams) {
	c.JSON(http.StatusOK, gin.H{})
}

// ServiceInstanceProvision provision a service instance
// (PUT /v2/service_instances/{instance_id})
func (s *server) ServiceInstanceProvision(c *gin.Context, _ string, _ ServiceInstanceProvisionParams) {
	body := ServiceInstanceProvisionJSONRequestBody{}
	err := c.BindJSON(&body)
	if err != nil {
		// TODO: handle error from AbortWithError
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse request body: %w", err))
		return
	}

	ctxParams := parseContext(body.Context)
	provision := s.client.NewProvision(
		ctxParams.OrgID,
		ctxParams.SpaceID,
		ctxParams.OrgName,
		ctxParams.SpaceName,
	)
	ctx := ctxutil.Ctx(c)
	if err = provision.ProvisionOrgSpacePolicy(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create org space policy: %w", err))
		return
	}
	if !ctx.IsEnableSpaceIdentity() {
		c.JSON(http.StatusCreated, gin.H{})
		return
	}
	if err = provision.ProvisionHostPolicy(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create space host policy: %w", err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// ServiceInstanceLastOperationGet get the last requested operation state for service instance
// (GET /v2/service_instances/{instance_id}/last_operation)
func (*server) ServiceInstanceLastOperationGet(c *gin.Context, _ string, _ ServiceInstanceLastOperationGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}
