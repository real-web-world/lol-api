package conf

import (
	"github.com/real-web-world/lol-api/pkg/logger"
)

type (
	AppConf struct {
		ProjectName string        `json:"projectName"`
		HTTPHost    string        `json:"host" default:"0.0.0.0" env:"httpHost"`
		HTTPPort    int           `json:"port" default:"8888"`
		Mode        string        `json:"mode" default:"release" env:"Mode"`
		Pyroscope   PyroscopeConf `json:"pyroscope"`
		Sentry      SentryConf    `json:"sentry"`
		Lang        LangConf      `json:"lang"`
		PProf       PProfConf     `json:"pprof"`
		Log         LogConf       `json:"log" required:"true"`
		Caller      string        `json:"caller" required:"true"`
		Redis       RedisConf     `json:"redis" required:"true"`
		Mysql       MysqlConf     `json:"mysql" required:"true"`
		Trace       TraceConf     `json:"trace" required:"true"`
		ProcDev     ProcDev       `json:"procDev"`
		Async       AsyncConf     `json:"asyncConf"`
	}

	PyroscopeConf struct {
		Enabled       bool   `json:"enabled" default:"true"`
		AppName       string `json:"appName"`
		ServerAddress string `json:"serverAddress"`
		SampleRate    uint32 `json:"sampleRate" default:"100"`
	}
	RedisItemConf struct {
		Host           string `required:"true" json:"host"`
		Port           int    `default:"6379" json:"port"`
		Pwd            string `default:"" json:"pwd"`
		DB             int    `default:"0" json:"db"`
		Pool           int    `default:"100" json:"pool"`
		CollectionName string `json:"collectionName"`
	}
	MysqlConf struct {
		Default MysqlItemConf `json:"default" required:"true"`
	}
	MysqlItemConf struct {
		Type        string `default:"mysql"`
		Host        string `required:"true"`
		Port        int    `required:"true"`
		UserName    string `required:"true"`
		Pwd         string `required:"true"`
		Charset     string `default:"utf8mb4"`
		Database    string `required:"true"`
		Prefix      string
		MaxIDleConn int `default:"10"`
		MaxOpenConn int `default:"100"`
		MaxLifeTime int `default:"60"`
	}
	RedisConf struct {
		Default RedisItemConf `json:"default" required:"true"`
	}
	AsyncConf struct {
		UpdateTicketSecond int `json:"updateTicketSecond" default:"5"`
	}
	TraceConf struct {
		Enabled    bool   `json:"enabled" default:"false"`
		ServerName string `json:"serverName" required:"true"`
		JaegerAddr string `json:"jaegerAddr" required:"true" env:"jaegerAddr"`
	}
	SentryConf struct {
		Enabled bool `json:"enabled" default:"true"`
		Dsn     string
	}
	LangConf struct {
		DefaultLang string `default:"zh" env:"defaultLang"`
	}
	PProfConf struct {
		Enabled bool   `default:"false" env:"enablePProf" json:"enabled"`
		Addr    string `json:"addr" default:":8889" env:"pprofAddr"`
	}
	LogConf struct {
		Level logger.LogLevelStr `json:"level" default:"info" env:"logLevel"`
	}
	ProcDev struct {
		AddDays int `json:"addDays" default:"0"`
	}
)
