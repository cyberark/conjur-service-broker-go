package servicebroker

import (
	"fmt"
	"net/http"

	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/gin-gonic/gin"
)

type hostBindServer struct {
	server
}

// ServiceBindingUnbinding deprovision a service binding
// (DELETE /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (s *hostBindServer) ServiceBindingUnbinding(c *gin.Context, instanceID string, bindingID string, params ServiceBindingUnbindingParams) {
	if err := validateServiceAndPlan(params.ServiceId, params.PlanId); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("validation failed: %w", err))
		return
	}

	// TODO: how to persist org and space id?
	bind := conjur.NewBind(s.client, "", "", bindingID)

	hostExists, err := bind.HostExists()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check host existance: %w", err))
		return
	}
	if !hostExists {
		c.AbortWithError(http.StatusGone, fmt.Errorf("host doesn't exists"))
		return
	}

	err = bind.DeletePolicy()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to delete policy: %w", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// ServiceBindingGet get a service binding
// (GET /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (*hostBindServer) ServiceBindingGet(c *gin.Context, instanceID string, bindingID string, params ServiceBindingGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}

// ServiceBindingBinding generate a service binding
// (PUT /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (s *hostBindServer) ServiceBindingBinding(c *gin.Context, instanceID string, bindingID string, params ServiceBindingBindingParams) {
	// TODO: Implement me
	body := ServiceBindingBindingJSONRequestBody{}
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

	ctxParams := parseContext(body.Context)

	bind := conjur.NewBind(s.client, ctxParams.OrgID, ctxParams.SpaceID, bindingID)
	hostExists, err := bind.HostExists()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check host existance: %w", err))
		return
	}
	if hostExists {
		c.AbortWithError(http.StatusConflict, fmt.Errorf("host already exists"))
		return
	}
	policy, err := bind.CreatePolicy()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create policy: %w", err))
		return
	}
	c.JSON(http.StatusCreated, ServiceBindingResponse{
		Credentials: object(policy),
	})
}

// ServiceBindingLastOperationGet get the last requested operation state for service binding
// (GET /v2/service_instances/{instance_id}/service_bindings/{binding_id}/last_operation)
func (*hostBindServer) ServiceBindingLastOperationGet(c *gin.Context, instanceID string, bindingID string, params ServiceBindingLastOperationGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}
