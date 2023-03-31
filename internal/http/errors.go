package http

import (
    "net/http"
    "strings"
    "unicode"

    "github.com/gin-gonic/gin"
    "golang.org/x/text/runes"
)

func errorsMiddleware(c *gin.Context) {
    c.Next()
    if len(c.Errors) > 0 {
        c.JSON(-1, gin.H{"error": strings.Join(c.Errors.Errors(), ", ")})
    }
    if c.IsAborted() && c.Writer.Size() == 0 { // request is aborted and body is missing
        statusCode := c.Writer.Status()
        if statusCode >= http.StatusBadRequest {
            c.JSON(-1, gin.H{"error": camelCasedStatus(statusCode)})
        }
    }
}

func camelCasedStatus(code int) string {
    var lastSpace bool
    return strings.Map(func(r rune) rune {
        if !runes.In(unicode.Letter).Contains(r) {
            lastSpace = true
            return -1
        }
        if lastSpace {
            lastSpace = false
            return unicode.ToUpper(r)
        }
        return unicode.ToLower(r)
    }, http.StatusText(code))
}
