package endpoints

import (
	"net/http"

	"github.com/nixys/nxs-go-appctx-example-restapi/ctx"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UsersList(cc *ctx.Ctx, c *gin.Context) RouteHandlerResponse {

	usrs, err := cc.User.GetAll()
	if err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("api users list")

		return RouteHandlerResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "internal error",
		}
	}

	return RouteHandlerResponse{
		StatusCode: http.StatusOK,
		RAWData:    usrs,
	}
}
