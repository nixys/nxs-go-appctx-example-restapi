package endpoints

import (
	"context"
	"net/http"
	"strings"

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

func Authorization(c context.Context, token string) gin.HandlerFunc {
	return func(gc *gin.Context) {

		at := gc.GetHeader("Authorization")
		if len(at) == 0 {
			gc.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		t := strings.TrimPrefix(at, "Bearer ")
		if t == at {
			gc.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if t == token {
			return
		}
		gc.AbortWithStatus(http.StatusUnauthorized)
	}
}
