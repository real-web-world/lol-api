package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	ginApp "github.com/real-web-world/bdk/gin"
	"github.com/real-web-world/lol-api/global"
	"github.com/real-web-world/lol-api/models"
	"github.com/real-web-world/lol-api/services/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func sgGetClientConf(ctx context.Context) (*models.Config, error) {
	res, err, _ := apiSg.Do(SgKeyClientCfg, func() (any, error) {
		db := global.DefaultDB
		m := models.NewCtxConfig(ctx, db)
		return m.GetClientConf()
	})
	return res.(*models.Config), err
}
func sgGetCurrVersion(ctx context.Context) ([]string, error) {
	res, err, _ := apiSg.Do(SgKeyCurrVersion, func() (any, error) {
		db := global.DefaultDB
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
			return nil, err
		}
		return []string{tag, downloadUrlPrefix}, nil
	})
	return res.([]string), err
}
func GetClientConf(c *gin.Context) {
	app := ginApp.GetApp(c)
	ctx := c.Request.Context()
	cfg, err := sgGetClientConf(ctx)
	if err != nil {
		logger.Error("获取客户端配置失败", zap.Error(err))
		app.ServerBad()
		return
	}
	app.Data(cfg.Val)
}

func GetCurrVersion(c *gin.Context) {
	app := ginApp.GetApp(c)
	ctx := c.Request.Context()
	res, err := sgGetCurrVersion(ctx)
	if err != nil {
		app.ServerBad()
		return
	}
	tag := res[0]
	downloadUrlPrefix := res[1]
	app.Data(gin.H{
		"versionTag":     tag,
		"downloadUrl":    fmt.Sprintf("%s/%s/hh-lol-prophet.exe", downloadUrlPrefix, tag),
		"zipDownloadUrl": fmt.Sprintf("%s/%s/hh-lol-prophet-%s.zip", downloadUrlPrefix, tag, tag),
	})
}
