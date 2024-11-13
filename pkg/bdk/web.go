package bdk

import (
	"net/http"
	"slices"
	"strings"

	"github.com/real-web-world/lol-api/global"
)

var (
	skipLogPathArr = []string{
		global.MetricsApi,
		global.VersionApi,
		global.StatusApi,
		global.FaviconReq,
	}
)

func IsSkipLogReq(req *http.Request, statusCode int) bool {
	isDevApi := strings.Index(req.RequestURI, global.DebugApiPrefix) == 0 ||
		strings.Index(req.RequestURI, global.SwaggerApiPrefix) == 0
	return isDevApi || statusCode == http.StatusNotFound || req.Method == http.MethodOptions ||
		slices.Index(skipLogPathArr, req.RequestURI) >= 0
}
