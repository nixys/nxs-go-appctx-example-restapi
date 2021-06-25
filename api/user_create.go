package api

import (
	"encoding/json"
	"net/http"

	"example/ctx"
	"example/modules/user"

	appctx "github.com/nixys/nxs-go-appctx/v2"

	"github.com/sirupsen/logrus"
)

type userCreateTx struct {
	User    userCreateDataTx `json:"user"`
	Message string           `json:"message"`
}

type userCreateDataTx struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type userCreateRx struct {
	Username string `json:"username"`
}

func userCreate(appCtx *appctx.AppContext) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rx := userCreateRx{}

		headersSet(w)

		if r.Method == http.MethodOptions {
			optionsResp(appCtx, w, r, http.StatusOK)
			return
		}

		cc := appCtx.CustomCtx().(*ctx.Ctx)

		// Check access for client
		if authCheck(cc.Conf.AuthKey, w, r) == false {
			userCreateResp(appCtx, w, r, http.StatusUnauthorized, user.Record{}, "auth check error")
			return
		}

		// Set request max body size
		r.Body = http.MaxBytesReader(w, r.Body, cc.Conf.ClientMaxBodySizeBytes)

		if r.ContentLength > cc.Conf.ClientMaxBodySizeBytes {
			userCreateResp(appCtx, w, r, http.StatusRequestEntityTooLarge, user.Record{}, "")
			return
		}

		// Retrieve json from body
		if err := json.NewDecoder(r.Body).Decode(&rx); err != nil {
			userCreateResp(appCtx, w, r, http.StatusBadRequest, user.Record{}, "incorrect body")
			return
		}

		if len(rx.Username) == 0 {
			userCreateResp(appCtx, w, r, http.StatusBadRequest, user.Record{}, "empty username")
			return
		}

		usr, err := user.Create(appCtx, user.CreateData{
			Username: rx.Username,
		})
		if err != nil {
			userCreateResp(appCtx, w, r, http.StatusInternalServerError, user.Record{}, err.Error())
			return
		}

		userCreateResp(appCtx, w, r, http.StatusOK, usr, "success")
		return
	})
}

func userCreateResp(appCtx *appctx.AppContext, w http.ResponseWriter, r *http.Request, httpCode int, usr user.Record, message string) {

	w.WriteHeader(httpCode)

	resp := userCreateTx{
		User: userCreateDataTx{
			ID:       usr.ID,
			Username: usr.Username,
			Password: usr.Password,
		},
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
