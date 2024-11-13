package cache

import (
	"github.com/real-web-world/lol-api/global"
)

var _ = fmtKey

func fmtKey(key string) string {
	return global.Conf.Redis.Default.CollectionName + ":" + key
}
