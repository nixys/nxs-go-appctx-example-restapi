package endpoints

import (
	"net/http"

	"github.com/nixys/nxs-go-appctx-example-restapi/ctx"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RouteHandlerResponse struct {
	StatusCode int
	RAWData    any
	Message    string
}

type RouteHandlers struct {
	Handler       RouteHandler
	DataTransform RouteDataTransformHandler
}

// Counter
type counterTx struct {
	Value int64 `json:"value"`
}

type RouteHandler func(*ctx.Ctx, *gin.Context) RouteHandlerResponse
type RouteDataTransformHandler func(any, string) any

func RouteHandlerDefault(cc *ctx.Ctx, handler RouteHandlers) gin.HandlerFunc {
	return func(c *gin.Context) {

		if handler.Handler == nil {
			cc.Log.Warn("route handler not specified")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		r := handler.Handler(cc, c)

		var d interface{}
		if handler.DataTransform != nil {
			d = handler.DataTransform(r.RAWData, r.Message)
		} else {
			d = r.RAWData
		}

		if d != nil {
			c.JSON(r.StatusCode, d)
		} else {
			c.String(r.StatusCode, r.Message)
		}
	}
}

func Logger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.WithFields(logrus.Fields{
			"type":      "accesslog",
			"remote":    c.RemoteIP(),
			"method":    c.Request.Method,
			"url":       c.Request.RequestURI,
			"code":      c.Writer.Status(),
			"userAgent": c.Request.UserAgent(),
		}).Info("request processed")
	}
}

func RequestSizeLimiter(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > limit {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

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
