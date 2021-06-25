package api

import (
	"encoding/json"
	"net/http"

	"example/ctx"

	"example/modules/user"

	appctx "github.com/nixys/nxs-go-appctx/v2"

	"github.com/sirupsen/logrus"
)

type usersGetTx struct {
	Users   []usersGetDataTx `json:"users"`
	Message string           `json:"message"`
}

type usersGetDataTx struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func usersGet(appCtx *appctx.AppContext) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		headersSet(w)

		if r.Method == http.MethodOptions {
			optionsResp(appCtx, w, r, http.StatusOK)
			return
		}

		cc := appCtx.CustomCtx().(*ctx.Ctx)

		// Check access for client
		if authCheck(cc.Conf.AuthKey, w, r) == false {
			usersGetResp(appCtx, w, r, http.StatusUnauthorized, []user.Record{}, "auth check error")
			return
		}

		usr, err := user.GetAll(appCtx)
		if err != nil {
			usersGetResp(appCtx, w, r, http.StatusInternalServerError, usr, err.Error())
			return
		}

		usersGetResp(appCtx, w, r, http.StatusOK, usr, "success")
		return
	})
}

func usersGetResp(appCtx *appctx.AppContext, w http.ResponseWriter, r *http.Request, httpCode int, users []user.Record, message string) {

	w.WriteHeader(httpCode)

	usrs := []usersGetDataTx{}

	for _, u := range users {
		usrs = append(usrs, usersGetDataTx{
			ID:       u.ID,
			Username: u.Username,
			Password: u.Password,
		})
	}

	resp := usersGetTx{
		Users:   usrs,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		appCtx.Log().WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"metod":       r.Method,
			"url":         r.URL,
			"code":        httpCode,
		}).Errorf("response send error: %v", err)
	}

	if httpCode != http.StatusOK {
		appCtx.Log().WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"metod":       r.Method,
			"url":         r.URL,
			"code":        httpCode,
			"message":     message,
		}).Warn("request processing finished with errors")
	}

	appCtx.Log().Infof("%s \"%s %s\" %d", r.RemoteAddr, r.Method, r.URL, httpCode)
}
