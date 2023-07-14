// Package servicebroker provides an implementation of the generated service broker server
package servicebroker

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.cyberng.com/Conjur-Enterprise/conjur-service-broker-go/internal/ctxutil"
	"github.cyberng.com/Conjur-Enterprise/conjur-service-broker-go/pkg/conjur"
)

// ServiceBindingUnbinding deprovision a service binding
// (DELETE /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (s *server) ServiceBindingUnbinding(c *gin.Context, _ string, bindingID string, _ ServiceBindingUnbindingParams) {
	if ctxutil.Ctx(c).IsEnableSpaceIdentity() {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	bind, err := s.client.FromBindingID(bindingID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check host existance: %w", err))
		return
	}
	hostExists, err := bind.HostExists()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check host existance: %w", err))
		return
	}
	if !hostExists {
		_ = c.AbortWithError(http.StatusGone, fmt.Errorf("host doesn't exists"))
		return
	}

	err = bind.DeleteBindHostPolicy()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to delete policy: %w", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// ServiceBindingGet get a service binding
// (GET /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (*server) ServiceBindingGet(c *gin.Context, _ string, _ string, _ ServiceBindingGetParams) {
	c.Status(http.StatusNotImplemented)
}

// ServiceBindingBinding generate a service binding
// (PUT /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (s *server) ServiceBindingBinding(c *gin.Context, _ string, bindingID string, _ ServiceBindingBindingParams) {
	body := ServiceBindingBindingJSONRequestBody{}
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse request body: %w", err))
		return
	}

	ctxParams := parseContext(body.Context)

	ctx := ctxutil.Ctx(c)
	bind := s.client.NewBind(ctxParams.OrgID, ctxParams.SpaceID, bindingID, ctx.IsEnableSpaceIdentity())
	hostExists, err := bind.HostExists()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check host existance: %w", err))
		return
	}
	if hostExists && !ctx.IsEnableSpaceIdentity() {
		_ = c.AbortWithError(http.StatusConflict, fmt.Errorf("host already exists"))
		return
	}
	if !hostExists && ctx.IsEnableSpaceIdentity() {
		_ = c.AbortWithError(http.StatusGone, fmt.Errorf("no space host identity found"))
		return
	}
	var policy *conjur.Policy
	if ctx.IsEnableSpaceIdentity() {
		policy, err = bind.BindSpacePolicy()
	} else {
		policy, err = bind.BindHostPolicy()
	}
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
	c.Status(http.StatusNotImplemented)
}
