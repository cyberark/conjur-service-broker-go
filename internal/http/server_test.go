//go:build !integration

package http

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.cyberng.com/Conjur-Enterprise/conjur-service-broker-go/internal/ctxutil"
	"go.uber.org/zap"
)

// func TestStartHTTPServer(t *testing.T) {
// 	tests := []struct {
// 		name string
// 	}{{}} // TODO: Add test cases.
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			StartHTTPServer()
// 		})
// 	}
// }

//
// func main() {
// 	var err error
// 	if len(os.Getenv("ERR")) > 0 {
// 		err = errors.New("error")
// 	}
// 	checkFatalErr(err)
// }
//
// func Test_checkFatalErr(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		err       bool
// 		wantFatal bool
// 	}{{
// 		"positive",
// 		false,
// 		false,
// 	}, {
// 		"negative",
// 		true,
// 		true,
// 	}}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Cleanup(cleanupEnv())
// 			if tt.err {
// 				err := os.Setenv("ERR", "")
// 				require.NoError(t, err)
// 			}
// 			err := exec.Command("go", "run", "./server_test.go").Run()
// 			println(err.Error())
// 		})
// 	}
// }

func Test_initCtx(t *testing.T) {
	type args struct {
		logger *zap.Logger
		cfg    *config
	}
	tests := []struct {
		name      string
		args      args
		wantEmpty bool
	}{{
		"not nil",
		args{
			logger: zap.NewNop(),
			cfg:    &config{},
		},
		false,
	}, {
		"nil",
		args{},
		true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := initCtx(tt.args.logger, tt.args.cfg)
			if tt.wantEmpty {
				require.Equal(t, ctxutil.NewContext(), got)
			} else {
				require.NotNil(t, got)
				require.NotEqualf(t, ctxutil.NewContext(), got, "expecting non empty context")
			}
		})
	}
}

//
// func Test_initLogger(t *testing.T) {
// 	type args struct {
// 		cfg *config
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		wantLogger  *zap.Logger
// 		wantCleanup func()
// 		wantErr     bool
// 	}{{}} // TODO: Add test cases.
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotLogger, _, err := initLogger(tt.args.cfg)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("initLogger() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotLogger, tt.wantLogger) {
// 				t.Errorf("initLogger() gotLogger = %v, want %v", gotLogger, tt.wantLogger)
// 			}
// 			// if !reflect.DeepEqual(gotCleanup, tt.wantCleanup) {
// 			// 	t.Errorf("initLogger() gotCleanup = %v, want %v", gotCleanup, tt.wantCleanup)
// 			// }
// 		})
// 	}
// }
//
// func Test_initServer(t *testing.T) {
// 	type args struct {
// 		cfg *config
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    servicebroker.ServerInterface
// 		wantErr bool
// 	}{{}} // TODO: Add test cases.
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := initServer(tt.args.cfg)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("initServer() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("initServer() got = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
