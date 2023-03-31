package conjur

//go:generate stringer -type=Privilege -linecomment -output priviledgestring.gen.go
//go:generate stringer -type=VariablePrivilege -linecomment -output variablepriviledgestring.gen.go

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
