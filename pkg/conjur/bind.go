package conjur

import (
	"io"
)

type bind struct {
	OrgID      string
	SpaceID    string
	InstanceID string
	BindID     string
	client     Client
}

type Bind interface {
	CreatePolicy() error
	Exists() (bool, error)
}

func NewBind(client Client, orgID, spaceID, instanceID, bindID string) Bind {
	res := &bind{
		OrgID:      orgID,
		SpaceID:    spaceID,
		InstanceID: instanceID,
		BindID:     bindID,
		client:     client,
	}
	return res
}

func (b *bind) CreatePolicy() error {
	//TODO implement me
	panic("implement me")
}

func (b *bind) Exists() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func createBindYAML(b *bind) (io.Reader, error) {
	policy := PolicyDocument{
		NewTag[Host](Host{
			Id:          b.BindID,
			Annotations: hostAnnotations(b),
		}),
		NewTag(Grant{
			Role:   NewRef[Layer](""),
			Member: NewRef[Layer](b.BindID),
		}),
	}
	return policyReader(policy)
}

func deleteBindYAML(b *bind) (io.Reader, error) {
	policy := PolicyDocument{
		NewTag(Delete{
			Record: NewRef[Host](b.BindID),
		}),
	}
	return policyReader(policy)
}

func hostAnnotations(b *bind) map[string]string {
	// TODO: support annotations
	return nil
}
