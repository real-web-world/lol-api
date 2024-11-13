package middleware

import (
	"github.com/gin-gonic/gin"

	ginApp "github.com/real-web-world/lol-api/pkg/gin"
)

func NotSaveResp(c *gin.Context) {
	app := ginApp.GetApp(c)
	app.SetNotSaveResp()
	c.Next()
}
