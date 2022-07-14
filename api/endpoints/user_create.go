package endpoints

import (
	"net/http"

	"example/modules/user"

	"github.com/gin-gonic/gin"
	appctx "github.com/nixys/nxs-go-appctx/v2"

	"github.com/sirupsen/logrus"
)

type userCreateTx struct {
	User    userDataTx `json:"user"`
	Message string     `json:"message"`
}

type userCreateRx struct {
	Username string `json:"username" binding:"required"`
}

func UserCreate(appCtx *appctx.AppContext, c *gin.Context) RouteHandlerResponse {

	rx := userCreateRx{}

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

	// Create user
	usr, err := user.Create(appCtx, user.CreateData{
		Username: rx.Username,
	})
	if err != nil {
		appCtx.Log().WithFields(logrus.Fields{
			"err": err,
		}).Warn("can't create user")
		return RouteHandlerResponse{
			StatusCode: http.StatusInternalServerError,
			RAWData:    user.Record{},
			Message:    "internal error",
		}
	}

	return RouteHandlerResponse{
		StatusCode: http.StatusOK,
		RAWData:    usr,
		Message:    "success",
	}
}

func UserCreateData(rd interface{}, message string) interface{} {
	usr := rd.(user.Record)
	return userCreateTx{
		User: userDataTx{
			ID:       usr.ID,
			Username: usr.Username,
			Password: usr.Password,
		},
		Message: message,
	}
}
