package conjur

import (
	"regexp"
	"strings"
)

//go:generate stringer -type=Kind -linecomment -output kind.gen.go

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
	k := Kind(-1)
	switch m[2] {
	case KindUser.String():
		k = KindUser
	case KindHost.String():
		k = KindHost
	case KindLayer.String():
		k = KindLayer
	case KindGroup.String():
		k = KindGroup
	case KindPolicy.String():
		k = KindPolicy
	case KindVariable.String():
		k = KindVariable
	case KindWebservice.String():
		k = KindWebservice
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
