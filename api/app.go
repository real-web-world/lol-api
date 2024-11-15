package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/real-web-world/bdk"
	ginApp "github.com/real-web-world/bdk/gin"

	apiProj "github.com/real-web-world/lol-api"
)

func ServerInfo(c *gin.Context) {
	app := ginApp.GetApp(c)
	now := time.Now()
	app.Data(bdk.ServerInfoData{
		Timestamp:   now.Unix(),
		TimestampMs: now.UnixMilli(),
	})
}
func DevHand(c *gin.Context) {
	app := ginApp.GetApp(c)
	d := &bdk.IDReq{}
	if err := c.ShouldBindJSON(d); err != nil {
		app.ValidError(err)
		return
	}
	time.Sleep(time.Millisecond)
	app.Data(gin.H{
		"id":   d.ID,
		"buff": 439612,
	})
}

func ShowVersion(c *gin.Context) {
	app := ginApp.GetApp(c)
	app.Data(bdk.VersionInfo{
		Version:   apiProj.APIVersion,
		Commit:    apiProj.Commit,
		BuildTime: apiProj.BuildTime,
		BuildUser: apiProj.BuildUser,
	})
}
func Status(c *gin.Context) {
	app := ginApp.GetApp(c)
	app.Data(gin.H{
		"mysql": true,
		"redis": true,
	})
}
