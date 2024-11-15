package routes

import (
	"github.com/gin-gonic/gin"
	bdkmid "github.com/real-web-world/bdk/gin/middleware"
	"github.com/real-web-world/lol-api/api"
)

func initAPIRoutes(r *gin.Engine) {
	initDevModule(r)
	initV1Module(r)
}

func initV1Module(r *gin.Engine) {
	lol := r.Group("lol")
	lol.POST("getCurrVersion", api.GetCurrVersion) // 获取当前版本和下载信息
	lolClient := lol.Group("client")
	lolClient.POST("getConf", api.GetClientConf) // 获取客户端配置
}

func initDevModule(r *gin.Engine) {
	dev := r.Group("", bdkmid.DevAccess)
	dev.POST("test", api.DevHand)
	dev.POST("serverInfo", api.ServerInfo)
}
