package api

import (
	"golang.org/x/sync/singleflight"
)

var (
	apiSg = singleflight.Group{}
)

const (
	SgKeyClientCfg   = "clientCfg"
	SgKeyCurrVersion = "currVersion"
)
