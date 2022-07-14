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

type userUpdateTx struct {
	User    userDataTx `json:"user"`
	Message string     `json:"message"`
}

type userUpdateRx struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
}

func UserUpdate(appCtx *appctx.AppContext, c *gin.Context) RouteHandlerResponse {

	rx := userUpdateRx{}

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

	// Fetch data from query
	if err := c.BindJSON(&rx); err != nil {
		appCtx.Log().WithFields(logrus.Fields{
			"err": err,
		}).Warn("can't parse body")
		return RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			RAWData:    user.Record{},
			Message:    "can't parse body",
		}
	}

	// Update user
	usr, err := user.Update(appCtx, user.UpdateData{
		ID:       id,
		Username: rx.Username,
		Password: rx.Password,
	})
	if err != nil {

		appCtx.Log().WithFields(logrus.Fields{
			"err": err,
		}).Warn("can't update user")

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

func UserUpdateData(rd interface{}, message string) interface{} {
	usr := rd.(user.Record)
	return userUpdateTx{
		User: userDataTx{
			ID:       usr.ID,
			Username: usr.Username,
			Password: usr.Password,
		},
		Message: message,
	}
}
