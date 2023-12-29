package endpoints

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/nixys/nxs-go-appctx-example-restapi/ctx"
	"github.com/nixys/nxs-go-appctx-example-restapi/misc"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UserDelete(c context.Context, cc *ctx.Ctx, gc *gin.Context) RouteHandlerResponse {

	id, err := strconv.ParseInt(gc.Param("id"), 10, 64)
	if err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("api user delete: can't parse id")

		return RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "incorrect user id",
		}
	}

	if err := cc.User.Delete(id); err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("api user delete")

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
	}
}
