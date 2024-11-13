package middleware

import (
	"time"

	sentryGin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"

	"github.com/real-web-world/lol-api/conf"
)

func GenerateSentryMiddleware(_ *conf.AppConf) gin.HandlerFunc {
	return sentryGin.New(sentryGin.Options{
		Repanic: true,
		Timeout: 3 * time.Second,
	})
}
