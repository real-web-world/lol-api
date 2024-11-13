package api

import (
	"time"

	"github.com/gin-gonic/gin"

	apiProj "github.com/real-web-world/lol-api"
	"github.com/real-web-world/lol-api/pkg/fastcurd"
	ginApp "github.com/real-web-world/lol-api/pkg/gin"
)

func ServerInfo(c *gin.Context) {
	app := ginApp.GetApp(c)
	now := time.Now()
	app.Data(ginApp.ServerInfo{
		Timestamp:   now.Unix(),
		TimestampMs: now.UnixMilli(),
	})
}

type (
	VersionInfo struct {
		Version   string `json:"version" example:"1.0.0"`
		Commit    string `json:"commit" example:"github commit sha256"`
		BuildTime string `json:"buildTime" example:"2006-01-02 15:04:05"`
		BuildUser string `json:"buildUser" example:"buffge"`
	}
	IDReq struct {
		ID uint64 `json:"id" binding:"required,min=1"`
	}
	IDArrReq struct {
		IDArr []uint64 `json:"idArr" binding:"required,min=1,max=20"`
	}
	ListCommonReq struct {
		LastID string `json:"lastID" binding:""`
		Page   uint   `json:"page" binding:"omitempty,min=0"`
		Limit  uint   `json:"limit" binding:"required,min=1,max=20"`
	}
	ListCommonResp struct {
		LastID string `json:"lastID"`
		More   bool   `json:"more"`
		Count  int64  `json:"count"`
	}
)

func DevHand(c *gin.Context) {
	app := ginApp.GetApp(c)
	d := &IDReq{}
	if err := c.ShouldBindJSON(d); err != nil {
		app.ValidError(err)
		return
	}
	app.Data(gin.H{
		"id":   d.ID,
		"buff": 439612,
	})
}

func ShowVersion(c *gin.Context) {
	app := ginApp.GetApp(c)
	json := &fastcurd.RetJSON{
		Code: fastcurd.CodeOk,
		Data: &VersionInfo{
			Version:   apiProj.APIVersion,
			Commit:    apiProj.Commit,
			BuildTime: apiProj.BuildTime,
			BuildUser: apiProj.BuildUser,
		},
	}
	app.JSON(json)
}
func Status(c *gin.Context) {
	app := ginApp.GetApp(c)
	app.Data(gin.H{
		"mysql": true,
		"redis": true,
	})
}
