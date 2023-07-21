package api

import (
	"example/api/endpoints"
	"example/ctx"

	"github.com/gin-gonic/gin"
	appctx "github.com/nixys/nxs-go-appctx/v2"
)

func RoutesSet(appCtx *appctx.AppContext) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(endpoints.Logger(appCtx))
	router.Use(endpoints.CORSMiddleware())

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	v1 := router.Group("/v1")
	{
		v1.Use(endpoints.RequestSizeLimiter(appCtx))
		v1.Use(endpoints.Authorize(cc.Conf.AuthKey))

		v1.GET("/users", endpoints.RouteHandlerDefault(appCtx, endpoints.RouteHandlers{
			Handler:       endpoints.UsersGet,
			DataTransform: endpoints.UsersGetData,
		}))

		v1.GET("/user/:id", endpoints.RouteHandlerDefault(appCtx, endpoints.RouteHandlers{
			Handler:       endpoints.UserGet,
			DataTransform: endpoints.UserGetData,
		}))

		v1.POST("/user", endpoints.RouteHandlerDefault(appCtx, endpoints.RouteHandlers{
			Handler:       endpoints.UserCreate,
			DataTransform: endpoints.UserCreateData,
		}))

		v1.PATCH("/user/:id", endpoints.RouteHandlerDefault(appCtx, endpoints.RouteHandlers{
			Handler:       endpoints.UserUpdate,
			DataTransform: endpoints.UserUpdateData,
		}))

		v1.DELETE("/user/:id", endpoints.RouteHandlerDefault(appCtx, endpoints.RouteHandlers{
			Handler:       endpoints.UserDelete,
			DataTransform: endpoints.UserDeleteData,
		}))
	}

	return router
}
