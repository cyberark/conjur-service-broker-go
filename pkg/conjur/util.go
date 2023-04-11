// Package conjur provides a wrapper around conjur go SDK
package conjur

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi"
)

//go:generate enumer -type=Kind -linecomment -text -output kind.gen.go

// Kind defines conjur resource kind
type Kind int

const (
	// KindUser One unique human
	KindUser Kind = iota // user
	// KindHost A single logical machine (in the broad sense, not just physical)
	KindHost // host
	// KindLayer A collection of hosts that have the same privileges
	KindLayer // layer
	// KindGroup A collection of users and groups that have the same privileges
	KindGroup // group
	// KindPolicy Privileges on policies enable a user to create and modify objects and permissions
	KindPolicy // policy
	// KindVariable A secret such as a password, API key, SSH key, etc
	KindVariable // variable
	// KindWebservice An HTTP(S) web service which performs sensitive operations
	KindWebservice // webservice
)

var conjurIDRegexp = regexp.MustCompile("^(?:(.*?)(?::|$))?(?:(.*?)(?::|$))?(.*?)$")

func parseID(id string) (account string, kind Kind, identifier string) {
	m := conjurIDRegexp.FindStringSubmatch(id)
	k, err := KindString(m[2])
	if err != nil {
		k = Kind(-1)
	}
	return m[1], k, m[3]
}

func composeID(account string, kind Kind, identifier string) string {
	var res []string
	if len(account) > 0 {
		res = append(res, account)
	}
	if len(kind.String()) > 0 {
		res = append(res, kind.String())
	}
	if len(identifier) > 0 {
		res = append(res, identifier)
	}
	return strings.Join(res, ":")
}

func apiKey(policy *conjurapi.PolicyResponse) (string, error) {
	if len(policy.CreatedRoles) < 1 {
		return "", fmt.Errorf("expecting at least one created role")
	}
	var roleID string
	var role conjurapi.CreatedRole
	for k, v := range policy.CreatedRoles {
		roleID = k
		role = v
		break
	}
	if roleID != role.ID {
		return "", fmt.Errorf("creatred role ID do not match %v != %v", roleID, role.ID)
	}
	return role.APIKey, nil
}
