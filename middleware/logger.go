package middleware

import (
	"net/http"
	"time"

	"github.com/real-web-world/lol-api/pkg/bdk"

	"github.com/gin-gonic/gin"

	"github.com/real-web-world/lol-api/global"
	ginApp "github.com/real-web-world/lol-api/pkg/gin"
)

func Logger(c *gin.Context) {
	c.Next()
	isSkipLog := bdk.IsSkipLogReq(c.Request, http.StatusOK)
	if isSkipLog {
		return
	}
	app := ginApp.GetApp(c)
	sqls := app.GetSqls()
	reqID := app.GetReqID()
	var totalTime time.Duration
	for _, sql := range sqls {
		execTime, _ := time.ParseDuration(sql.ExecTime)
		totalTime += execTime
	}
	if len(sqls) > 0 {
		global.Logger.Infow("sql",
			"reqID", reqID,
			"sqlCount", len(sqls),
			"totalTime", totalTime.String(),
			"sqls", sqls,
		)
	}
}
