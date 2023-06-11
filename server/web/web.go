package web

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/inarithefox/partsy/server/app"
	"github.com/inarithefox/partsy/server/public/logger"
	"github.com/inarithefox/partsy/server/public/model"
	"github.com/inarithefox/partsy/server/public/utils"
)

type Web struct {
	srv    *app.Server
	Router *mux.Router
}

func New(srv *app.Server) *Web {
	web := &Web{
		srv:    srv,
		Router: srv.Router,
	}

	return web
}

func HandleNotFound(a *app.App, w http.ResponseWriter, r *http.Request) {
	err := model.NewAppError("NotFound", "App Error", nil, "", http.StatusNotFound)
	ipAddress := utils.GetIPAddress(r)
	logger.Debug(fmt.Sprintf("%s %s %s %d %s", ipAddress, r.URL.Path, r.Method, http.StatusNotFound, r.Proto))

	if IsApiCall(a, r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.StatusCode)
		err.DetailedError = fmt.Sprintf("There does not appear to be an API call for the url '%s'", r.URL.Path)
		w.Write([]byte(err.ToJSON()))
	} else {
		http.NotFound(w, r)
	}
}

func IsApiCall(a *app.App, r *http.Request) bool {
	subpath, _ := utils.GetSubpathFromConfig(a.Config())

	return strings.HasPrefix(r.URL.Path, path.Join(subpath, "api")+"/")
}
