package conjur

import "testing"

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
