package ctx

import (
	units "github.com/docker/go-units"
	conf "github.com/nixys/nxs-go-conf"
)

type confOpts struct {
	Bind              string    `conf:"bind" conf_extraopts:"default=0.0.0.0:8080"`
	LogFile           string    `conf:"logfile" conf_extraopts:"default=stdout"`
	LogLevel          string    `conf:"loglevel" conf_extraopts:"default=info"`
	PidFile           string    `conf:"pidfile"`
	ClientMaxBodySize string    `conf:"clientMaxBodySize" conf_extraopts:"default=36m"`
	TLS               tlsConf   `conf:"tls"`
	AuthKey           string    `conf:"authKey" conf_extraopts:"required"`
	MySQL             mysqlConf `conf:"mysql" conf_extraopts:"required"`
	CounterInterval   int64     `conf:"counterInterval" conf_extraopts:"default=5"`

	ClientMaxBodySizeBytes int64
}

type tlsConf struct {
	CertFile string `conf:"certfile"`
	KeyFie   string `conf:"keyfile"`
}

type mysqlConf struct {
	Host     string `conf:"host" conf_extraopts:"required"`
	DB       string `conf:"db" conf_extraopts:"required"`
	User     string `conf:"user" conf_extraopts:"required"`
	Password string `conf:"password" conf_extraopts:"required"`
}

func confRead(confPath string) (confOpts, error) {

	var c confOpts

	err := conf.Load(&c, conf.Settings{
		ConfPath:    confPath,
		ConfType:    conf.ConfigTypeYAML,
		UnknownDeny: true,
	})
	if err != nil {
		return c, err
	}

	c.ClientMaxBodySizeBytes, err = units.RAMInBytes(c.ClientMaxBodySize)
	if err != nil {
		return c, err
	}

	return c, err
}
