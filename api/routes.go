package api

import (
	"context"
	"time"

	"github.com/nixys/nxs-go-appctx-example-restapi/api/endpoints"
	"github.com/nixys/nxs-go-appctx-example-restapi/ctx"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RoutesSet(c context.Context, cc *ctx.Ctx) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(endpoints.Logger(cc.Log))

	router.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     cc.API.CORS.AllowOrigins,
				AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding"},
				ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)

	v1 := router.Group("/v1")
	{
		v1.Use(endpoints.RequestSizeLimiter(cc.API.ClientMaxBodySizeBytes))
		v1.Use(endpoints.Authorization(c, cc.API.AuthToken))

		user := v1.Group("/user")
		{
			user.GET("", endpoints.RouteHandlerDefault(c, cc, endpoints.RouteHandlers{
				Handler:       endpoints.UsersList,
				DataTransform: endpoints.UsersData,
			}))

			user.GET("/:id", endpoints.RouteHandlerDefault(c, cc, endpoints.RouteHandlers{
				Handler:       endpoints.UserGet,
				DataTransform: endpoints.UserData,
			}))

			user.POST("", endpoints.RouteHandlerDefault(c, cc, endpoints.RouteHandlers{
				Handler:       endpoints.UserCreate,
				DataTransform: endpoints.UserData,
			}))

			user.PATCH("/:id", endpoints.RouteHandlerDefault(c, cc, endpoints.RouteHandlers{
				Handler:       endpoints.UserUpdate,
				DataTransform: endpoints.UserData,
			}))

			user.DELETE("/:id", endpoints.RouteHandlerDefault(c, cc, endpoints.RouteHandlers{
				Handler: endpoints.UserDelete,
			}))
		}
	}

	return router
}
