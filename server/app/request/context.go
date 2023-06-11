package request

import (
	"context"

	"github.com/inarithefox/partsy/server/public/model"
)

type Context struct {
	session        model.Session
	requestId      string
	ipAddress      string
	path           string
	userAgent      string
	acceptLanguage string

	context context.Context
}

func NewContext(ctx context.Context, requestId, ipAddress, path, userAgent, acceptLanguage string, session model.Session) *Context {
	return &Context{
		session:        session,
		requestId:      requestId,
		ipAddress:      ipAddress,
		path:           path,
		userAgent:      userAgent,
		acceptLanguage: acceptLanguage,
		context:        ctx,
	}
}

func (c *Context) Session() *model.Session {
	return &c.session
}

func (c *Context) SetAcceptLanguage(s string) {
	c.acceptLanguage = s
}

func (c *Context) SetContext(ctx context.Context) {
	c.context = ctx
}

func (c *Context) SetPath(s string) {
	c.path = s
}

func (c *Context) SetRequestId(s string) {
	c.requestId = s
}

func (c *Context) SetIPAddress(s string) {
	c.ipAddress = s
}

func (c *Context) SetUserAgent(s string) {
	c.userAgent = s
}
