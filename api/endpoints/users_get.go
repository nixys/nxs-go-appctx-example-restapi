package endpoints

import (
	"net/http"

	"example/modules/user"

	"github.com/gin-gonic/gin"
	appctx "github.com/nixys/nxs-go-appctx/v2"
	"github.com/sirupsen/logrus"
)

type usersGetTx struct {
	Users   []userDataTx `json:"users"`
	Message string       `json:"message"`
}

func UsersGet(appCtx *appctx.AppContext, c *gin.Context) RouteHandlerResponse {

	// Get users
	usrs, err := user.GetAll(appCtx)
	if err != nil {
		appCtx.Log().WithFields(logrus.Fields{
			"err": err,
		}).Warn("can't get all users")
		return RouteHandlerResponse{
			StatusCode: http.StatusInternalServerError,
			RAWData:    []user.Record{},
			Message:    "internal error",
		}
	}

	return RouteHandlerResponse{
		StatusCode: http.StatusOK,
		RAWData:    usrs,
		Message:    "success",
	}
}

func UsersGetData(rd interface{}, message string) interface{} {

	users := rd.([]user.Record)

	usrs := []userDataTx{}

	for _, u := range users {
		usrs = append(usrs, userDataTx{
			ID:       u.ID,
			Username: u.Username,
			Password: u.Password,
		})
	}

	return usersGetTx{
		Users:   usrs,
		Message: message,
	}
}
