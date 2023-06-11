package web

import (
	"strings"

	"github.com/inarithefox/partsy/server/app"
	"github.com/inarithefox/partsy/server/app/request"
	"github.com/inarithefox/partsy/server/public/model"
)

type Context struct {
	App           *app.App
	AppContext    *request.Context
	Params        *Params
	Err           *model.AppError
	siteUrlHeader string
}

func (c *Context) GetSiteURLHeader() string {
	return c.siteUrlHeader
}

func (c *Context) SetSiteURLHeader(url string) {
	c.siteUrlHeader = strings.TrimRight(url, "/")
}
