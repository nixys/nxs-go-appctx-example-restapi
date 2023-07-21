package endpoints

import (
	"example/misc"
	"example/modules/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	appctx "github.com/nixys/nxs-go-appctx/v2"
	"github.com/sirupsen/logrus"
)

type userDeleteTx struct {
	Message string `json:"message"`
}

func UserDelete(appCtx *appctx.AppContext, c *gin.Context) RouteHandlerResponse {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		appCtx.Log().WithFields(logrus.Fields{
			"err": err,
		}).Warn("can't parse id")
		return RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "incorrect id",
		}
	}

	// Delete user
	if err := user.Delete(appCtx, id); err != nil {

		appCtx.Log().WithFields(logrus.Fields{
			"err": err,
		}).Warn("can't delete user")

		switch err {
		case misc.ErrNotFound:
			return RouteHandlerResponse{
				StatusCode: http.StatusNotFound,
				Message:    "user not found",
			}
		default:
			return RouteHandlerResponse{
				StatusCode: http.StatusInternalServerError,
				RAWData:    user.Record{},
				Message:    "internal error",
			}
		}
	}

	return RouteHandlerResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
	}
}

func UserDeleteData(rd interface{}, message string) interface{} {
	return userDeleteTx{
		Message: message,
	}
}
