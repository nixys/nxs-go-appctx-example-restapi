package endpoints

import (
	"net/http"
	"strconv"

	"example/misc"
	"example/modules/user"

	"github.com/gin-gonic/gin"
	appctx "github.com/nixys/nxs-go-appctx/v2"
	"github.com/sirupsen/logrus"
)

type userGetTx struct {
	User    userDataTx `json:"user"`
	Message string     `json:"message"`
}

func UserGet(appCtx *appctx.AppContext, c *gin.Context) RouteHandlerResponse {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		appCtx.Log().WithFields(logrus.Fields{
			"err": err,
		}).Warn("can't parse id")
		return RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			RAWData:    user.Record{},
			Message:    "incorrect id",
		}
	}

	// Get users
	usr, err := user.Get(appCtx, id)
	if err != nil {

		appCtx.Log().WithFields(logrus.Fields{
			"err": err,
		}).Warn("can't get user")

		switch err {
		case misc.ErrNotFound:
			return RouteHandlerResponse{
				StatusCode: http.StatusNotFound,
				RAWData:    user.Record{},
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
		RAWData:    usr,
		Message:    "success",
	}
}

func UserGetData(rd interface{}, message string) interface{} {

	user := rd.(user.Record)

	return userGetTx{
		User: userDataTx{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
		},
		Message: message,
	}
}
