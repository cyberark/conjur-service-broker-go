package conjur

type PolicyDocument []*Tag[Policy]

type Policy struct {
	Id          string            `yaml:"id"`
	Owner       Ref[User]         `yaml:"owner,omitempty"`
	Body        []interface{}     `yaml:"body,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type Group struct {
	Id    string    `yaml:"id"`
	Owner Ref[User] `yaml:"owner"`
}

type User struct {
	Id    string    `yaml:"id"`
	Owner Ref[User] `yaml:"owner"`
}

type Layer string

type Grant struct {
	Role   *Tag[any] `yaml:"role"`
	Member *Tag[any] `yaml:"member"`
}
