package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func errorsMiddleware(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		c.JSON(-1, gin.H{"error": strings.Join(c.Errors.Errors(), ", ")})
	}
}
