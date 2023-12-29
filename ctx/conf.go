package ctx

import (
	"fmt"

	conf "github.com/nixys/nxs-go-conf"
)

type confOpts struct {
	LogFile  string `conf:"logfile" conf_extraopts:"default=stdout"`
	LogLevel string `conf:"loglevel" conf_extraopts:"default=info"`

	API   apiConf   `conf:"api" conf_extraopts:"required"`
	MySQL mysqlConf `conf:"mysql" conf_extraopts:"required"`
}

type apiConf struct {
	Bind              string   `conf:"bind" conf_extraopts:"default=0.0.0.0:8080"`
	TLS               *tlsConf `conf:"tls"`
	ClientMaxBodySize string   `conf:"clientMaxBodySize" conf_extraopts:"default=36m"`
	AuthToken         string   `conf:"authToken" conf_extraopts:"required"`
	CORS              corsConf `conf:"cors" conf_extraopts:"required"`
}

type mysqlConf struct {
	Host     string `conf:"host" conf_extraopts:"default=127.0.0.1"`
	Port     int    `conf:"port" conf_extraopts:"default=3306"`
	DB       string `conf:"db" conf_extraopts:"required"`
	User     string `conf:"user" conf_extraopts:"required"`
	Password string `conf:"password" conf_extraopts:"required"`
}

type tlsConf struct {
	CertFile string `conf:"certfile" conf_extraopts:"required"`
	KeyFie   string `conf:"keyfile" conf_extraopts:"required"`
}

type corsConf struct {
	AllowOrigins []string `conf:"allowOrigins" conf_extraopts:"required"`
}

func confRead(confPath string) (confOpts, error) {

	var c confOpts

	err := conf.Load(&c, conf.Settings{
		ConfPath:    confPath,
		ConfType:    conf.ConfigTypeYAML,
		UnknownDeny: true,
	})
	if err != nil {
		return c, fmt.Errorf("conf init: %w", err)
	}

	return c, nil
}
