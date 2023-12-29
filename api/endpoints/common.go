package endpoints

import (
	"context"
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

type RouteHandler func(context.Context, *ctx.Ctx, *gin.Context) RouteHandlerResponse
type RouteDataTransformHandler func(context.Context, any, string) any

func RouteHandlerDefault(c context.Context, cc *ctx.Ctx, handler RouteHandlers) gin.HandlerFunc {
	return func(gc *gin.Context) {

		if handler.Handler == nil {
			cc.Log.Warn("route handler not specified")
			gc.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		r := handler.Handler(c, cc, gc)

		var d interface{}
		if handler.DataTransform != nil {
			d = handler.DataTransform(c, r.RAWData, r.Message)
		} else {
			d = r.RAWData
		}

		if d != nil {
			gc.JSON(r.StatusCode, d)
		} else {
			gc.String(r.StatusCode, r.Message)
		}
	}
}

func Logger(log *logrus.Logger) gin.HandlerFunc {
	return func(gc *gin.Context) {
		gc.Next()
		log.WithFields(logrus.Fields{
			"type":      "accesslog",
			"remote":    gc.RemoteIP(),
			"method":    gc.Request.Method,
			"url":       gc.Request.RequestURI,
			"code":      gc.Writer.Status(),
			"userAgent": gc.Request.UserAgent(),
		}).Info("request processed")
	}
}

func RequestSizeLimiter(limit int64) gin.HandlerFunc {
	return func(gc *gin.Context) {
		if gc.Request.ContentLength > limit {
			gc.AbortWithStatus(http.StatusRequestEntityTooLarge)
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(gc *gin.Context) {

		gc.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		gc.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		gc.Writer.Header().Set("Access-Control-Allow-Headers", "X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, X-Auth-Health-Key, X-Auth-Key, If-Modified-Since, Cache-Control, Content-Type")
		gc.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")

		if gc.Request.Method == http.MethodOptions {
			gc.AbortWithStatus(http.StatusNoContent)
			return
		}
	}
}

func Authorize(c context.Context, authKey string) gin.HandlerFunc {
	return func(gc *gin.Context) {
		if gc.GetHeader("X-Auth-Key") == authKey {
			return
		}
		gc.AbortWithStatus(http.StatusUnauthorized)
	}
}
