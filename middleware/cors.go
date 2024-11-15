package middleware

import (
	"regexp"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/real-web-world/lol-api/conf"
	"github.com/real-web-world/lol-api/global"
)

var (
	allowOriginRegex = regexp.MustCompile(".+?\\.buffge\\.com(:\\d+)?$")
)

func Cors(_ *conf.AppConf) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if global.IsDevMode() {
				return true
			}
			return allowOriginRegex.MatchString(origin)
		},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT"},
		AllowHeaders:     []string{"content-type", "x-requested-with", "token", "locale", "Bf-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
