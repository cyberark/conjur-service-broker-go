// Package conjur provides a wrapper around conjur go SDK
package conjur

import (
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/stretchr/testify/require"
)

func Test_parseID(t *testing.T) {
	tests := []struct {
		name           string
		args           string
		wantAccount    string
		wantKind       Kind
		wantIdentifier string
	}{{
		"full",
		"abc:host:ghi",
		"abc",
		KindHost,
		"ghi",
	}, {
		"account and kind",
		"abc:user",
		"abc",
		KindUser,
		"",
	}, {
		"just account",
		"abc",
		"abc",
		Kind(-1),
		"",
	}, {
		"wit additional colon in id",
		"abc:group:ghi:123",
		"abc",
		KindGroup,
		"ghi:123",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAccount, gotKind, gotIdentifier := parseID(tt.args)
			if gotAccount != tt.wantAccount {
				t.Errorf("parseID() gotAccount = %v, want %v", gotAccount, tt.wantAccount)
			}
			if gotKind != tt.wantKind {
				t.Errorf("parseID() gotKind = %v, want %v", gotKind, tt.wantKind)
			}
			if gotIdentifier != tt.wantIdentifier {
				t.Errorf("parseID() gotIdentifier = %v, want %v", gotIdentifier, tt.wantIdentifier)
			}
		})
	}
}

func Test_composeID(t *testing.T) {
	type args struct {
		account    string
		kind       Kind
		identifier string
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		"full",
		args{
			account:    "account",
			kind:       KindHost,
			identifier: "id",
		},
		"account:host:id",
	}, {
		"no account",
		args{
			account:    "",
			kind:       KindHost,
			identifier: "id",
		},
		"host:id",
	}, {
		"no kind",
		args{
			account:    "account",
			kind:       -1,
			identifier: "id",
		},
		"account:id",
	}, {
		"no id",
		args{
			account:    "account",
			kind:       KindPolicy,
			identifier: "",
		},
		"account:policy",
	}, {
		"just account",
		args{
			account:    "account",
			kind:       -1,
			identifier: "",
		},
		"account",
	}, {
		"empty",
		args{
			account:    "",
			kind:       -1,
			identifier: "",
		},
		"",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := composeID(tt.args.account, tt.args.kind, tt.args.identifier)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_apiKey(t *testing.T) {
	tests := []struct {
		name    string
		policy  *conjurapi.PolicyResponse
		want    string
		wantErr bool
	}{{
		"positive",
		&conjurapi.PolicyResponse{
			CreatedRoles: map[string]conjurapi.CreatedRole{"role": {
				APIKey: "api-key",
				ID:     "role",
			}},
		},
		"api-key",
		false,
	}, {
		"mismatched id",
		&conjurapi.PolicyResponse{
			CreatedRoles: map[string]conjurapi.CreatedRole{"role": {
				APIKey: "api-key",
				ID:     "id",
			}},
		},
		"",
		true,
	}, {
		"empty",
		nil,
		"",
		true,
	}, {
		"no role",
		&conjurapi.PolicyResponse{},
		"",
		true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := apiKey(tt.policy)
			if tt.wantErr {
				require.Error(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_slashJoin(t *testing.T) {
	tests := []struct {
		name  string
		elems []string
		want  string
	}{{
		"simple",
		[]string{"a", "b", "c"},
		"a/b/c",
	}, {
		"single",
		[]string{"abc"},
		"abc",
	}, {
		"single with slashes",
		[]string{"/abc/"},
		"/abc/",
	}, {
		"empty",
		[]string{},
		"",
	}, {
		"nil",
		nil,
		"",
	}, {
		"with slash",
		[]string{"/a/b/", "/c"},
		"/a/b/c",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slashJoin(tt.elems...); got != tt.want {
				t.Errorf("slashJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}
