// Package conjur provides a wrapper around conjur go SDK
package conjur

import (
	"bytes"
	"io"

	"github.com/doodlesbykumbi/conjur-policy-go/pkg/conjurpolicy"
	"gopkg.in/yaml.v3"
)

func policyReader(policy conjurpolicy.PolicyStatements) (io.Reader, error) {
	res := new(bytes.Buffer)
	if len(policy) == 0 {
		res.WriteString("\n")
		return res, nil
	}
	encoder := yaml.NewEncoder(res)
	err := encoder.Encode(policy)
	if err != nil {
		return nil, err
	}
	return res, err
}
