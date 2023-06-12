package web

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/inarithefox/partsy/server/app"
	"github.com/inarithefox/partsy/server/app/request"
	"github.com/inarithefox/partsy/server/public/logger"
	"github.com/inarithefox/partsy/server/public/model"
	"github.com/inarithefox/partsy/server/public/utils"
)

type Handler struct {
	Srv            *app.Server
	HandleFunc     func(*Context, http.ResponseWriter, *http.Request)
	HandlerName    string
	RequireSession bool
	RequireMfa     bool
	IsStatic       bool
}

func (w *Web) NewHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &Handler{
		Srv:            w.srv,
		HandleFunc:     h,
		HandlerName:    GetHandlerName(h),
		RequireSession: false,
		RequireMfa:     false,
		IsStatic:       false,
	}
}

func GetHandlerName(h func(*Context, http.ResponseWriter, *http.Request)) string {
	handlerName := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	pos := strings.LastIndex(handlerName, ".")
	if pos != -1 && len(handlerName) > pos {
		handlerName = handlerName[pos+1:]
	}

	return handlerName
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w = NewWrappedWriter(w)
	appInstance := app.New(app.ServerConnector(h.Srv))

	requestId := model.NewId()
	var statusCode string

	ipAddress := utils.GetIPAddress(r)

	defer func() {
		logger.Debug(fmt.Sprintf("%s %s %s %s %s", ipAddress, r.URL.Path, r.Method, statusCode, r.Proto))
	}()

	c := &Context{
		AppContext: &request.Context{},
		App:        appInstance,
	}

	c.AppContext.SetRequestId(requestId)
	c.AppContext.SetIPAddress(ipAddress)
	c.AppContext.SetPath(r.URL.Path)
	c.AppContext.SetContext(context.Background())
	c.Params = ParamsFromRequest(r)

	subpath, err := utils.GetSubpathFromConfig(c.App.Config())
	if err != nil {
		logger.Error(err, "unable to get subpath from config")
	}
	siteUrlHeader := fmt.Sprintf("%s://%s%s", utils.GetProtocol(r), r.Host, subpath)
	c.SetSiteURLHeader(siteUrlHeader)

	w.Header().Set(model.HeaderRequestId, requestId)
	w.Header().Set(model.HeaderVersionId, model.CurrentVersion)

	w.Header().Set("Referrer-Policy", "no-referrer")

	if c.Err == nil {
		h.HandleFunc(c, w, r)
	}

	statusCode = strconv.Itoa(w.(*ResponseWriterWrapper).StatusCode())
}
