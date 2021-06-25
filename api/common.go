package api

import (
	"net/http"

	appctx "github.com/nixys/nxs-go-appctx/v2"
)

// headersSet sets default headers
func headersSet(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, X-Auth-Health-Key, X-Auth-Key, If-Modified-Since, Cache-Control, Content-Type")
}

// optionsResp sends default response to any OPTIONS request
func optionsResp(appCtx *appctx.AppContext, w http.ResponseWriter, r *http.Request, httpCode int) {
	w.WriteHeader(httpCode)
	appCtx.Log().Infof("%s \"%s %s\" %d", r.RemoteAddr, r.Method, r.URL, httpCode)
}

func authCheck(authKey string, w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("X-Auth-Key") == authKey {
		return true
	}
	return false
}
