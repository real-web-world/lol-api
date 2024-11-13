package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/real-web-world/lol-api/api"
	mid "github.com/real-web-world/lol-api/middleware"
)

func initAPIRoutes(r *gin.Engine) {
	initDevModule(r)
	initV1Module(r)
}

func initV1Module(r *gin.Engine) {
	lol := r.Group("lol")
	lolClient := lol.Group("client")
	lolClient.POST("getConf", api.GetClientConf)
}

func initDevModule(r *gin.Engine) {
	dev := r.Group("", mid.DevAccess)
	dev.POST("test", api.DevHand)
	dev.POST("serverInfo", api.ServerInfo)
}
