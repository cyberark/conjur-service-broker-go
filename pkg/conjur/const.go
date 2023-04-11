// Package conjur provides a wrapper around conjur go SDK
package conjur

//go:generate enumer -type=Privilege -linecomment -output priviledgestring.gen.go
//go:generate enumer -type=VariablePrivilege -linecomment -output variablepriviledgestring.gen.go

// Privilege defines an enum describing possible values for policy, user, host, group, layer
type Privilege int

const (
	// PrivilegeRead read
	PrivilegeRead Privilege = iota // read
	// PrivilegeUpdate update
	PrivilegeUpdate // update
	// PrivilegeCreate create
	PrivilegeCreate // create
)

// VariablePrivilege defines an enum describing possible values for variable
type VariablePrivilege int

const (
	// VariablePrivilegeRead read
	VariablePrivilegeRead VariablePrivilege = iota // read
	// VariablePrivilegeExecute execute
	VariablePrivilegeExecute // execute
	// VariablePrivilegeUpdate update
	VariablePrivilegeUpdate // update
)

const spaceHostAPIKey = "space-host-api-key" // #nosec G101
