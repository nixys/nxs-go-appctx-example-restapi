package api

import (
	"net/http"

	appctx "github.com/nixys/nxs-go-appctx/v2"

	"github.com/gorilla/mux"
)

// RoutesSet sets handlers for API server
func RoutesSet(appCtx *appctx.AppContext) *mux.Router {

	r := mux.NewRouter()

	r.Handle("/users", usersGet(appCtx)).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/user", userCreate(appCtx)).Methods(http.MethodPut, http.MethodOptions)

	return r
}
