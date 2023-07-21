package endpoints

import (
	"example/ctx"
	"net/http"

	"github.com/gin-gonic/gin"
	appctx "github.com/nixys/nxs-go-appctx/v2"
	"github.com/sirupsen/logrus"
)

type userDataTx struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RouteHandlerResponse struct {
	StatusCode int
	RAWData    interface{}
	Message    string
}

type RouteHandlers struct {
	Handler       RouteHandler
	DataTransform RouteDataTransformHandler
}

type RouteHandler func(*appctx.AppContext, *gin.Context) RouteHandlerResponse
type RouteDataTransformHandler func(interface{}, string) interface{}

func RouteHandlerDefault(appCtx *appctx.AppContext, handlers RouteHandlers) gin.HandlerFunc {
	return func(c *gin.Context) {

		if handlers.Handler == nil {
			appCtx.Log().Warn("route handler not specified")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		r := handlers.Handler(appCtx, c)

		var d interface{}
		if handlers.DataTransform != nil {
			d = handlers.DataTransform(r.RAWData, r.Message)
		} else {
			d = r.RAWData
		}

		c.JSON(r.StatusCode, d)
	}
}

func Logger(appCtx *appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		appCtx.Log().WithFields(logrus.Fields{
			"type":      "accesslog",
			"remote":    c.RemoteIP(),
			"method":    c.Request.Method,
			"url":       c.Request.RequestURI,
			"code":      c.Writer.Status(),
			"userAgent": c.Request.UserAgent(),
		}).Info("request processed")
	}
}

func RequestSizeLimiter(appCtx *appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		cc := appCtx.CustomCtx().(*ctx.Ctx)
		if c.Request.ContentLength > cc.Conf.ClientMaxBodySizeBytes {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, X-Auth-Health-Key, X-Auth-Key, If-Modified-Since, Cache-Control, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
	}
}

func Authorize(authKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Auth-Key") == authKey {
			return
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
