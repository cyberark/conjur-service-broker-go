// Package http implements the http communication layer of the conjur service broker
package http

import (
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/runes"
)

func errorsMiddleware(c *gin.Context) {
	c.Header("Content-Type", "application/json") // this is needed to avoid basic auth enforcing content type text/plain https://github.com/gin-gonic/gin/issues/1453
	c.Next()
	if c.Writer.Size() > 0 { // body was already sent
		return
	}
	if len(c.Errors) > 0 {
		c.JSON(-1, gin.H{"error": camelCasedStatus(c.Writer.Status()), "description": strings.Join(c.Errors.Errors(), ", ")})
		return
	}
	if c.IsAborted() {
		statusCode := c.Writer.Status()
		if statusCode >= http.StatusBadRequest {
			c.JSON(-1, gin.H{"error": camelCasedStatus(statusCode)})
		}
	}
}

func camelCasedStatus(code int) string {
	lastSpace := true
	return strings.Map(func(r rune) rune {
		switch {
		case !runes.In(unicode.Letter).Contains(r):
			lastSpace = true
			return -1
		case lastSpace:
			lastSpace = false
			return unicode.ToUpper(r)
		default:
			return unicode.ToLower(r)
		}
	}, http.StatusText(code))
}
