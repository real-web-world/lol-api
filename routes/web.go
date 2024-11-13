package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/real-web-world/lol-api/api"
	"github.com/real-web-world/lol-api/global"
	mid "github.com/real-web-world/lol-api/middleware"
)

func initWebRoutes(r *gin.Engine) {
	r.GET(global.VersionApi, mid.NotSaveResp, mid.NotTrace, api.ShowVersion)
	r.GET(global.StatusApi, mid.NotSaveResp, mid.NotTrace, api.Status)
	dev := r.Group("", mid.DevAccess)
	dev.GET(global.MetricsApi, mid.NotSaveResp, gin.WrapH(promhttp.Handler()))
}
