package http

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/caarlos0/env/v7"
	"github.com/cyberark/conjur-service-broker/pkg/conjur"
)

func Test_newConfig(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		want    map[string]interface{}
		wantErr error
	}{{
		"missing required",
		map[string]string{},
		nil,
		env.EmptyEnvVarError{},
	}, {
		"required set but empty value",
		parseEnvs("CONJUR_ACCOUNT=;CONJUR_APPLIANCE_URL=http://localhost:8082;CONJUR_AUTHN_API_KEY=api-key;CONJUR_AUTHN_LOGIN=host/service-broker;CONJUR_POLICY=cf;DEBUG=true;ENABLE_SPACE_IDENTITY=true;SECURITY_USER_NAME=test;SECURITY_USER_PASSWORD=test"),
		nil,
		env.EmptyEnvVarError{},
	}, {
		"positive",
		parseEnvs("CONJUR_ACCOUNT=dev;CONJUR_APPLIANCE_URL=http://localhost:8082;CONJUR_AUTHN_API_KEY=api-key;CONJUR_AUTHN_LOGIN=host/service-broker;CONJUR_POLICY=cf;DEBUG=true;ENABLE_SPACE_IDENTITY=true;SECURITY_USER_NAME=test;SECURITY_USER_PASSWORD=test"),
		map[string]interface{}{"CONJUR_VERSION": uint32(5)},
		nil,
	}}
	for _, tt := range tests {
		t.Cleanup(cleanupEnv())
		t.Run(tt.name, func(t *testing.T) {
			initEnvs(t, tt.envs)
			got, err := newConfig()
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("newConfig() error = %v", err)
				return
			}
			if tt.want == nil && got != nil {
				t.Errorf("newConfig() expected nil, got %v", got)
				return
			}
			if len(tt.want) > 0 && !assertConfigContains(got, tt.want) {
				t.Errorf("newConfig() expected config to contain %v, got %+v", tt.want, got)
				return
			}
		})
	}
}

func assertConfigContains(got *config, want map[string]interface{}) bool {
	cfg := configMap(*got, got.Config)
	for k, v := range want {
		cfgVal, ok := cfg[k]
		if !ok || cfgVal != v {
			return false
		}
	}
	return true
}

func configMap(cfg ...interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for _, c := range cfg {
		v := reflect.ValueOf(c)
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			tag, ok := t.Field(i).Tag.Lookup("env")
			if !ok {
				continue
			}
			envName := strings.Split(tag, ",")[0]
			res[envName] = v.Field(i).Interface()
		}
	}
	return res
}

func parseEnvs(envs string) map[string]string {
	res := map[string]string{}
	for _, e := range strings.Split(envs, ";") {
		v := strings.SplitN(e, "=", 2)
		res[v[0]] = v[1]
	}
	return res
}

func initEnvs(t *testing.T, envs map[string]string) {
	for k, v := range envs {
		t.Setenv(k, v)
	}
}

func Test_validate(t *testing.T) {
	type args struct {
		cfg config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{{
		"validate positive",
		args{cfg: config{
			Config: conjur.Config{
				ConjurVersion: 5,
			},
		}},
		false,
	}, {
		"validate negative",
		args{cfg: config{
			Config: conjur.Config{
				ConjurVersion: 4,
			},
		}},
		true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
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
