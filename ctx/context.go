package ctx

import (
	"fmt"
	"os"
	"time"

	"github.com/docker/go-units"
	"github.com/nixys/nxs-go-appctx-example-restapi/ds/primedb"
	"github.com/nixys/nxs-go-appctx-example-restapi/modules/user"

	appctx "github.com/nixys/nxs-go-appctx/v3"
	"github.com/sirupsen/logrus"
)

// Ctx defines application custom context
type Ctx struct {
	CounterInterval time.Duration
	Log             *logrus.Logger
	User            *user.User
	API             apiCtx
}

type apiCtx struct {
	Bind                   string
	TLS                    *apiTLSCtx
	ClientMaxBodySizeBytes int64
	AuthToken              string
	CORS                   corsCtx
}

type apiTLSCtx struct {
	CertFile string
	KeyFie   string
}

type corsCtx struct {
	AllowOrigins []string
}

func AppCtxInit() (any, error) {

	c := &Ctx{}

	args, err := ArgsRead()
	if err != nil {
		return nil, err
	}

	conf, err := confRead(args.ConfigPath)
	if err != nil {
		// Write to temp logger
		tmpLogError("ctx init: %s", err.Error())
		return nil, err
	}

	c.Log, err = logInit(conf.LogFile, conf.LogLevel)
	if err != nil {
		tmpLogError("ctx init: %s", err.Error())
		return nil, err
	}

	primeDB, err := primedb.Connect(primedb.Settings{
		Host:     conf.MySQL.Host,
		Port:     conf.MySQL.Port,
		Database: conf.MySQL.DB,
		User:     conf.MySQL.User,
		Password: conf.MySQL.Password,
	})
	if err != nil {
		c.Log.Errorf("ctx init: %s", err.Error())
		return nil, err
	}

	c.User = user.Init(user.InitSettings{
		DB: primeDB,
	})

	bts, err := units.RAMInBytes(conf.API.ClientMaxBodySize)
	if err != nil {
		c.Log.Errorf("ctx init: parse client max body size: %s", err.Error())
		return nil, err
	}

	c.API = apiCtx{
		Bind: conf.API.Bind,
		TLS: func() *apiTLSCtx {
			if conf.API.TLS == nil {
				return nil
			}
			return &apiTLSCtx{
				CertFile: conf.API.TLS.CertFile,
				KeyFie:   conf.API.TLS.KeyFie,
			}
		}(),
		ClientMaxBodySizeBytes: bts,
		AuthToken:              conf.API.AuthToken,
		CORS: corsCtx{
			AllowOrigins: func() []string {
				orgns := []string{}
				for _, o := range conf.API.CORS.AllowOrigins {
					orgns = append(orgns, o)
				}
				return orgns
			}(),
		},
	}

	return c, nil
}

func tmpLogError(format string, args ...interface{}) {
	l, _ := appctx.DefaultLogInit(os.Stderr, logrus.InfoLevel, &logrus.JSONFormatter{})
	l.Errorf(format, args...)
}

func logInit(file, level string) (*logrus.Logger, error) {

	var (
		f   *os.File
		err error
	)

	switch file {
	case "stdout":
		f = os.Stdout
	case "stderr":
		f = os.Stderr
	default:
		f, err = os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			return nil, fmt.Errorf("log init: %w", err)
		}
	}

	// Validate log level
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("log init: %w", err)
	}

	return appctx.DefaultLogInit(f, l, &logrus.JSONFormatter{})
}
