package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	ginApp "github.com/real-web-world/bdk/gin"
	"github.com/real-web-world/lol-api/global"
	"github.com/real-web-world/lol-api/models"
	"github.com/real-web-world/lol-api/services/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

func GetCurrVersion(c *gin.Context) {
	app := ginApp.GetApp(c)
	db := global.DefaultDB
	ctx := c.Request.Context()
	m := models.NewCtxConfig(ctx, db)
	var tag, downloadUrlPrefix string
	g := errgroup.Group{}
	g.Go(func() error {
		cfg, err := m.GetCurrVersion()
		if err != nil {
			logger.Error("获取当前版本失败", zap.Error(err))
			return err
		}
		tag = cfg.Val
		return nil
	})
	g.Go(func() error {
		cfg, err := m.GetDownloadUrlPrefix()
		if err != nil {
			logger.Error("获取下载地址失败", zap.Error(err))
			return err
		}
		downloadUrlPrefix = cfg.Val
		return nil
	})
	if err := g.Wait(); err != nil {
		app.ServerBad()
		return
	}
	app.Data(gin.H{
		"versionTag":     tag,
		"downloadUrl":    fmt.Sprintf("%s/%s/hh-lol-prophet.exe", downloadUrlPrefix, tag),
		"zipDownloadUrl": fmt.Sprintf("%s/%s/hh-lol-prophet-%s.zip", downloadUrlPrefix, tag, tag),
	})
}
