package conjur

import (
	"bytes"
	"io"

	"github.com/doodlesbykumbi/conjur-policy-go/pkg/conjurpolicy"
	"gopkg.in/yaml.v3"
)

func policyReader(policy conjurpolicy.PolicyStatements) (io.Reader, error) {
	res := new(bytes.Buffer)
	encoder := yaml.NewEncoder(res)
	err := encoder.Encode(policy)
	if err != nil {
		return nil, err
	}
	return res, err
}
