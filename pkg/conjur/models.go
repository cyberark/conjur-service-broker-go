package conjur

type PolicyDocument []Tag

type Policy struct {
	Id          string            `yaml:"id"`
	Owner       Ref               `yaml:"owner,omitempty"`
	Body        []Tag             `yaml:"body,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type Group struct {
	Id    string `yaml:"id"`
	Owner Ref    `yaml:"owner"`
}

type User struct {
	Id    string `yaml:"id"`
	Owner Ref    `yaml:"owner"`
}

type Layer string

type Grant struct {
	Role   Ref `yaml:"role"`
	Member Ref `yaml:"member"`
}

type Host struct {
	Id          string            `yaml:"id"`
	Owner       Ref               `yaml:"owner,omitempty"`
	Body        []Tag             `yaml:"body,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type Delete struct {
	Record Ref `yaml:"record"`
}
