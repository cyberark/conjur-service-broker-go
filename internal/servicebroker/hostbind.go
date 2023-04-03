package servicebroker

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type hostBindServer struct {
	server
}

// ServiceBindingUnbinding deprovision a service binding
// (DELETE /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (*hostBindServer) ServiceBindingUnbinding(c *gin.Context, instanceID string, bindingID string, params ServiceBindingUnbindingParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}

// ServiceBindingGet get a service binding
// (GET /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (*hostBindServer) ServiceBindingGet(c *gin.Context, instanceID string, bindingID string, params ServiceBindingGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}

// ServiceBindingBinding generate a service binding
// (PUT /v2/service_instances/{instance_id}/service_bindings/{binding_id})
func (*hostBindServer) ServiceBindingBinding(c *gin.Context, instanceID string, bindingID string, params ServiceBindingBindingParams) {
	// TODO: Implement me

	c.Status(http.StatusNotImplemented)
}

// ServiceBindingLastOperationGet get the last requested operation state for service binding
// (GET /v2/service_instances/{instance_id}/service_bindings/{binding_id}/last_operation)
func (*hostBindServer) ServiceBindingLastOperationGet(c *gin.Context, instanceID string, bindingID string, params ServiceBindingLastOperationGetParams) {
	// TODO: Implement me
	c.Status(http.StatusNotImplemented)
}
