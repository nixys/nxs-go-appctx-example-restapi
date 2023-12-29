package endpoints

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/nixys/nxs-go-appctx-example-restapi/ctx"
	"github.com/nixys/nxs-go-appctx-example-restapi/misc"
	"github.com/nixys/nxs-go-appctx-example-restapi/modules/user"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

type userUpdateRx struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func UserUpdate(c context.Context, cc *ctx.Ctx, gc *gin.Context) RouteHandlerResponse {

	rx := userUpdateRx{}

	id, err := strconv.ParseInt(gc.Param("id"), 10, 64)
	if err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("api user update: can't parse id")

		return RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "incorrect user id",
		}
	}

	// Fetch data from query
	if err := gc.BindJSON(&rx); err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("api user update: can't parse body")

		return RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "can't parse body",
		}
	}

	usr, err := cc.User.Update(user.UpdateData{
		ID:       id,
		Username: rx.Username,
		Password: rx.Password,
	})
	if err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("api user update")

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
