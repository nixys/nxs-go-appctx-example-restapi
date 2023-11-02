package endpoints

import (
	"net/http"

	"github.com/nixys/nxs-go-appctx-example-restapi/ctx"
	"github.com/nixys/nxs-go-appctx-example-restapi/modules/user"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

type userCreateRx struct {
	Username string `json:"username" binding:"required"`
}

func UserCreate(cc *ctx.Ctx, c *gin.Context) RouteHandlerResponse {

	rx := userCreateRx{}

	// Fetch data from query
	if err := c.BindJSON(&rx); err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("api user create: can't parse body")

		return RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "can't parse body",
		}
	}

	usr, err := cc.User.Create(user.CreateData{
		Username: rx.Username,
	})
	if err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("api user create")

		return RouteHandlerResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "internal error",
		}
	}

	return RouteHandlerResponse{
		StatusCode: http.StatusOK,
		RAWData:    usr,
	}
}
