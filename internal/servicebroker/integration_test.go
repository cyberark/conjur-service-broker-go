package servicebroker

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cyberark/conjur-service-broker/internal/ctxutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func ginTestCtx(t *testing.T, method, url string, body string, enableSpaceIdentity bool) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ctx := ctxutil.NewContext()
	ctx = ctx.WithEnableSpaceIdentity(enableSpaceIdentity)
	c.Set("service-broker-context", ctx)
	var b io.Reader
	if len(body) > 0 {
		b = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, b)
	require.NoError(t, err)
	c.Request = req
	return w, c
}
