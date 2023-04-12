// Package servicebroker provides an implementation of the generated service broker server
package servicebroker

import (
	"fmt"
	"net/http"

	"github.com/cyberark/conjur-service-broker/internal/ctxutil"
	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/gin-gonic/gin"
)

// ServiceBindingUnbinding deprovision a service binding
// (DELETE /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (s *server) ServiceBindingUnbinding(c *gin.Context, _ string, bindingID string, _ ServiceBindingUnbindingParams) {
	bind, err := conjur.FromBindingID(s.client, bindingID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check host existance: %w", err))
		return
	}
	hostExists, err := bind.HostExists(ctxutil.Ctx(c))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check host existance: %w", err))
		return
	}
	if !hostExists {
		_ = c.AbortWithError(http.StatusGone, fmt.Errorf("host doesn't exists"))
		return
	}

	err = bind.DeletePolicy(ctxutil.Ctx(c))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to delete policy: %w", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// ServiceBindingGet get a service binding
// (GET /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (*server) ServiceBindingGet(c *gin.Context, _ string, _ string, _ ServiceBindingGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}

// ServiceBindingBinding generate a service binding
// (PUT /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (s *server) ServiceBindingBinding(c *gin.Context, _ string, bindingID string, _ ServiceBindingBindingParams) {
	// TODO: Implement me
	body := ServiceBindingBindingJSONRequestBody{}
	err := c.BindJSON(&body)
	if err != nil {
		// TODO: handle error from AbortWithError
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse request body: %w", err))
		return
	}

	ctxParams := parseContext(body.Context)

	bind := conjur.NewBind(s.client, ctxParams.OrgID, ctxParams.SpaceID, bindingID)
	hostExists, err := bind.HostExists(ctxutil.Ctx(c))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check host existance: %w", err))
		return
	}
	if hostExists {
		_ = c.AbortWithError(http.StatusConflict, fmt.Errorf("host already exists"))
		return
	}
	policy, err := bind.CreatePolicy(ctxutil.Ctx(c))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create policy: %w", err))
		return
	}
	c.JSON(http.StatusCreated, ServiceBindingResponse{
		Credentials: object(policy),
	})
}

// ServiceBindingLastOperationGet get the last requested operation state for service binding
// (GET /v2/service_instances/{instance_id}/service_bindings/{binding_id}/last_operation)
func (*server) ServiceBindingLastOperationGet(c *gin.Context, _ string, _ string, _ ServiceBindingLastOperationGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}
