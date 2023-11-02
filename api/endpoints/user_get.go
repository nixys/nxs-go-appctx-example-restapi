package endpoints

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nixys/nxs-go-appctx-example-restapi/ctx"
	"github.com/nixys/nxs-go-appctx-example-restapi/misc"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UserGet(cc *ctx.Ctx, c *gin.Context) RouteHandlerResponse {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("api user get: can't parse id")

		return RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "incorrect user id",
		}
	}

	usr, err := cc.User.Get(id)
	if err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("api user get")

		switch {
		case errors.Is(err, misc.ErrNotFound):
			return RouteHandlerResponse{
				StatusCode: http.StatusNotFound,
				Message:    "user not found",
			}
		default:
			return RouteHandlerResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal error",
			}
		}
	}

	return RouteHandlerResponse{
		StatusCode: http.StatusOK,
		RAWData:    usr,
	}
}
