package conf

type (
	AppConf struct {
		ProjectName string        `json:"projectName" required:"true"`
		HTTPHost    string        `json:"host" default:"0.0.0.0" env:"httpHost"`
		HTTPPort    int           `json:"port" default:"8888"`
		Mode        string        `json:"mode" default:"release" env:"Mode"`
		Pyroscope   PyroscopeConf `json:"pyroscope"`
		Lang        LangConf      `json:"lang"`
		PProf       PProfConf     `json:"pprof"`
		Log         LogConf       `json:"log" required:"true"`
		Redis       RedisConf     `json:"redis" required:"true"`
		Mysql       MysqlConf     `json:"mysql" required:"true"`
		Otel        OtelConf      `json:"otel" required:"true"`
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
		Host               string `json:"host" required:"true"`
		Port               int    `json:"port" required:"true"`
		UserName           string `json:"userName" required:"true"`
		Pwd                string `json:"pwd" required:"true"`
		Charset            string `json:"charset" default:"utf8mb4"`
		Database           string `json:"database" required:"true"`
		Prefix             string `json:"prefix"`
		MaxIDleConn        int    `json:"maxIDleConn" default:"10"`
		MaxOpenConn        int    `json:"maxOpenConn" default:"100"`
		MaxLifeTimeMinutes int    `json:"maxLifeTimeMinutes" default:"60"`
	}
	RedisConf struct {
		Default RedisItemConf `json:"default" required:"true"`
	}
	OtelConf struct {
		Enabled  bool   `json:"enabled" default:"false"`
		Endpoint string `json:"endpoint" required:"true"`
	}
	LangConf struct {
		DefaultLang string `default:"zh" env:"defaultLang"`
	}
	PProfConf struct {
		Enabled bool `default:"false" json:"enabled"`
	}
	LogConf struct {
		Level string `json:"level" default:"info" env:"logLevel"`
	}
)
