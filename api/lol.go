package api

import (
	"github.com/gin-gonic/gin"
	"github.com/real-web-world/lol-api/global"
	"github.com/real-web-world/lol-api/models"
	ginApp "github.com/real-web-world/lol-api/pkg/gin"
	"github.com/real-web-world/lol-api/services/logger"
	"go.uber.org/zap"
)

func GetClientConf(c *gin.Context) {
	app := ginApp.GetApp(c)
	db := global.DefaultDB
	ctx := c.Request.Context()
	m := models.NewCtxConfig(ctx, db)
	cfg, err := m.GetClientConf()
	if err != nil {
		logger.Error("获取客户端配置失败", zap.Error(err))
		app.ServerBad()
		return
	}
	app.Data(cfg.Val)
}
