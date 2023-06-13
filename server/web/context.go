package web

import (
	"net/http"
	"strings"

	"github.com/inarithefox/partsy/server/app"
	"github.com/inarithefox/partsy/server/app/request"
	"github.com/inarithefox/partsy/server/public/model"
)

type Context struct {
	App                  *app.App
	AppContext           *request.Context
	Params               *Params
	Err                  *model.AppError
	SiteUrlHeader        string
	GraphQLOperationName string
}

func (c *Context) GetSiteURLHeader() string {
	return c.SiteUrlHeader
}

func (c *Context) RequirePartId() *Context {
	if c.Err != nil {
		return c
	}

	if !model.IsValidId(c.Params.PartId) {
		c.SetInvalidURLParameter("part_id")
	}

	return c
}

func (c *Context) SetInvalidURLParameter(parameter string) {
	c.Err = model.NewAppError("Context", "api.context.invalid_url_parameter", map[string]any{"Name": parameter}, "", http.StatusBadRequest)
}

func (c *Context) SetInvalidParameterWithErr(parameter string, er error) {
	c.Err = model.NewAppError("Context", "api.context.invalid_parameter", map[string]any{"Name": parameter}, "", http.StatusBadRequest)
}

func (c *Context) SetSiteURLHeader(url string) {
	c.SiteUrlHeader = strings.TrimRight(url, "/")
}
