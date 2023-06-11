package api

import (
	"net/http"

	"github.com/inarithefox/partsy/server/web"
)

type handlerFunc func(*web.Context, http.ResponseWriter, *http.Request)

func (api *Api) APIHandler(h handlerFunc) http.Handler {
	handler := &web.Handler{
		Srv:            api.srv,
		HandleFunc:     h,
		HandlerName:    web.GetHandlerName(h),
		RequireSession: false,
		RequireMfa:     false,
		IsStatic:       false,
	}

	return handler
}

func (api *Api) APISessionRequired(h handlerFunc) http.Handler {
	handler := &web.Handler{
		Srv:            api.srv,
		HandleFunc:     h,
		HandlerName:    web.GetHandlerName(h),
		RequireSession: true,
		RequireMfa:     true,
		IsStatic:       false,
	}

	return handler
}

func (api *Api) APISessionRequiredWithoutMFA(h handlerFunc) http.Handler {
	handler := &web.Handler{
		Srv:            api.srv,
		HandleFunc:     h,
		HandlerName:    web.GetHandlerName(h),
		RequireSession: true,
		RequireMfa:     false,
		IsStatic:       false,
	}

	return handler
}
