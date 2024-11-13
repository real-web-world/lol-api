package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	remoteIPKey = "x-remote-ip"
)

func InjectRemoteIP(c *gin.Context) {
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), remoteIPKey, c.ClientIP()))
	c.Next()
}
