package api

import (
	"github.com/nixys/nxs-go-appctx-example-restapi/api/endpoints"
	"github.com/nixys/nxs-go-appctx-example-restapi/ctx"

	"github.com/gin-gonic/gin"
)

func RoutesSet(cc *ctx.Ctx) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(endpoints.Logger(cc.Log))
	router.Use(endpoints.CORSMiddleware())

	v1 := router.Group("/v1")
	{
		v1.Use(endpoints.RequestSizeLimiter(cc.API.ClientMaxBodySizeBytes))
		v1.Use(endpoints.Authorize(cc.Conf.API.AuthToken))

		user := v1.Group("/user")
		{
			user.GET("", endpoints.RouteHandlerDefault(cc, endpoints.RouteHandlers{
				Handler:       endpoints.UsersList,
				DataTransform: endpoints.UsersData,
			}))

			user.GET("/:id", endpoints.RouteHandlerDefault(cc, endpoints.RouteHandlers{
				Handler:       endpoints.UserGet,
				DataTransform: endpoints.UserData,
			}))

			user.POST("", endpoints.RouteHandlerDefault(cc, endpoints.RouteHandlers{
				Handler:       endpoints.UserCreate,
				DataTransform: endpoints.UserData,
			}))

			user.PATCH("/:id", endpoints.RouteHandlerDefault(cc, endpoints.RouteHandlers{
				Handler:       endpoints.UserUpdate,
				DataTransform: endpoints.UserData,
			}))

			user.DELETE("/:id", endpoints.RouteHandlerDefault(cc, endpoints.RouteHandlers{
				Handler: endpoints.UserDelete,
			}))
		}
	}

	return router
}
