package conjur

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cyberark/conjur-service-broker-go/pkg/conjur/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConfig_NewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{{
		"no follower",
		Config{
			ConjurAccount:      "account",
			ConjurApplianceURL: "https://conjur.local",
		},
		func(at assert.TestingT, got interface{}, args ...interface{}) bool {
			t := at.(*testing.T)
			require.NotNil(t, got)
			c, ok := got.(*client)
			require.True(t, ok)
			return c.client.GetConfig().ApplianceURL == c.roClient.GetConfig().ApplianceURL
		},
		assert.NoError,
	}, {
		"with follower",
		Config{
			ConjurAccount:      "account",
			ConjurApplianceURL: "https://conjur.local",
			ConjurFollowerURL:  "https://follower.local",
		},
		func(at assert.TestingT, got interface{}, args ...interface{}) bool {
			t := at.(*testing.T)
			require.NotNil(t, got)
			c, ok := got.(*client)
			require.True(t, ok)
			return c.client.GetConfig().ApplianceURL != c.roClient.GetConfig().ApplianceURL
		},
		assert.NoError,
	}, {
		"missing appliance url",
		Config{
			ConjurAccount: "account",
		},
		func(at assert.TestingT, got interface{}, args ...interface{}) bool {
			return got == nil
		},
		assert.Error,
	}}
	for _, tt := range tests {
		t.Cleanup(cleanupEnv())
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.NewClient(nil)
			tt.wantErr(t, err)
			tt.want(t, got)
		})
	}
}

func cleanupEnv() func() {
	values := os.Environ()
	os.Clearenv()
	return func() {
		os.Clearenv()
		for _, v := range values {
			parts := strings.SplitN(v, "=", 2)
			_ = os.Setenv(parts[0], parts[1])
		}
	}
}

func Test_client_NewBind(t *testing.T) {
	tests := []struct {
		name                string
		enableSpaceIdentity bool
		wantHostID          string
	}{{
		"host identity",
		false,
		"account:host:/orgID/spaceID/bindingID",
	}, {
		"space identity",
		true,
		"account:host:/orgID/spaceID",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{config: &Config{ConjurAccount: "account"}}
			b := c.NewBind("orgID", "spaceID", "bindingID", tt.enableSpaceIdentity)
			got, ok := b.(*bind)
			require.True(t, ok)
			assert.Equal(t, tt.wantHostID, got.hostID)
		})
	}
}

func Test_client_NewProvision(t *testing.T) {
	testStr := "test_str"
	type args struct {
		orgID     string
		spaceID   string
		orgName   *string
		spaceName *string
	}
	tests := []struct {
		name string
		args args
		want args
	}{{
		"without names",
		args{
			orgID:   "orgID",
			spaceID: "spaceID",
		},
		args{
			orgID:   "orgID",
			spaceID: "spaceID",
		},
	}, {
		"with names",
		args{
			orgName:   &testStr,
			spaceName: &testStr,
		},
		args{
			orgName:   &testStr,
			spaceName: &testStr,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{config: &Config{ConjurAccount: "account"}}
			b := c.NewProvision(tt.args.orgID, tt.args.spaceID, tt.args.orgName, tt.args.spaceName)
			got, ok := b.(*provision)
			require.True(t, ok)
			assert.Equal(t, tt.want.orgID, got.orgID)
			assert.Equal(t, tt.want.spaceID, got.spaceID)
			if tt.want.orgName == nil {
				assert.Empty(t, got.orgName)
			} else {
				assert.Equal(t, tt.want.orgName, &got.orgName)
			}
			if tt.want.spaceName == nil {
				assert.Empty(t, got.spaceName)
			} else {
				assert.Equal(t, tt.want.spaceName, &got.spaceName)
			}
		})
	}
}

func Test_client_platformAnnotation(t *testing.T) {
	tests := []struct {
		name       string
		conjurResp interface{}
		want       string
		wantErr    assert.ErrorAssertionFunc
	}{{
		"with platform annotations",
		map[string]interface{}{"annotations": []map[string]string{{"name": "platform", "value": "pivotalcloudfoundry"}}},
		"pivotalcloudfoundry",
		assert.NoError,
	}, {
		"no platform annotation",
		map[string]interface{}{"annotations": []map[string]string{{"name": "not-platform", "value": "pivotalcloudfoundry"}}},
		"",
		assert.NoError,
	}, {
		"invalid json for unmarshal",
		map[string]interface{}{"annotations": map[string]string{"name": "not-platform", "value": "pivotalcloudfoundry"}},
		"",
		assert.Error,
	}, {
		"invalid json for marshal",
		map[string]interface{}{"annotations": func() {}},
		"",
		assert.Error,
	}, {
		"error from resource",
		errors.New("error"),
		"",
		assert.Error,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mocks.NewMockClient(t)
			var conjurErr error
			var conjurResp interface{}
			if e, isErr := tt.conjurResp.(error); isErr {
				conjurErr = e
			} else {
				conjurResp = tt.conjurResp
			}
			c.On("Resource", mock.Anything).Return(conjurResp, conjurErr).Once()
			client := client{roClient: c, config: &Config{}}
			got, err := client.platformAnnotation()
			tt.wantErr(t, err)
			assert.Equalf(t, tt.want, got, "platformAnnotation()")
			c.AssertExpectations(t)
		})
	}
}
