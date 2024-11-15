package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/real-web-world/bdk"
	bdkmid "github.com/real-web-world/bdk/gin/middleware"
	"github.com/real-web-world/lol-api/api"
)

func initWebRoutes(r *gin.Engine) {
	r.GET(bdk.VersionApi, bdkmid.NotTrace, api.ShowVersion)
	r.GET(bdk.StatusApi, bdkmid.NotTrace, api.Status)
	dev := r.Group("", bdkmid.DevAccess)
	dev.GET(bdk.MetricsApi, bdkmid.NotTrace, gin.WrapH(promhttp.Handler()))
}
